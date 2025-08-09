package handler

import (
    "net/http"
    "github.com/locne/auth-service/internal/interface/repository"
    "github.com/locne/auth-service/internal/usecase"
    "github.com/gin-gonic/gin"
    "fmt"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  []string    `json:"errors,omitempty"`
}

type LoginDto struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type CreateUserDto struct {
    Username string `json:"username"`
	Email string `json:"email"`
    Password string `json:"password"`
	Elo string `json:"elo"`
}

func Login(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req LoginDto
        if err := c.ShouldBindJSON(&req); err != nil {
            fmt.Println("BindJSON error:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        userInfo, accessToken, refreshToken, err := usecase.Login(userRepo, req.Username, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }

		c.SetCookie(
			"access_token",
			accessToken,
			7*24*60*60, // maxAge: 7 day
			"/",        // path
			"",         // domain 
			false,      // secure: false (true if HTTPS)
			true,       // httpOnly: true
		)
		c.SetCookie(
			"refresh_token",
			refreshToken,
			7*24*60*60,
			"/",
			"",
			false,
			true,
		)
		c.JSON(http.StatusOK, APIResponse{
			Status:  "success",
			Message: "Login successful",
			Data:    userInfo,
		})
    }
}


func RegisterAuthRoutes(router *gin.Engine, userRepo repository.UserRepository ) {
    api := router.Group("/api/v1/auth")
    {
        api.POST("/login", Login(userRepo))
    }
}