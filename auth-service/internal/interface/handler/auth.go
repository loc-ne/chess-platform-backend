package handler

import (
    "net/http"
    "github.com/locne/auth-service/internal/interface/repository"
    "github.com/locne/auth-service/internal/usecase"
    "github.com/locne/auth-service/internal/entity"
    "github.com/locne/auth-service/internal/infrastructure/messagebroker"
    "github.com/rabbitmq/amqp091-go"
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
	Elo int `json:"elo"`
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

func Register(userRepo repository.UserRepository, ch *amqp091.Channel) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req CreateUserDto
        if err := c.ShouldBindJSON(&req); err != nil {
            fmt.Println("BindJSON error:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        user := entity.User{
        Email:    req.Email,
        Username: req.Username,
        Password: req.Password,
        }

        userInfo, err := usecase.Register(userRepo, user)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }

        if ch != nil {
            if err := messagebroker.PublishPlayerRegister(ch, userInfo.ID, req.Elo); err != nil {
                fmt.Println("Publish player register error:", err)
            }
        }

		c.JSON(http.StatusOK, APIResponse{
			Status:  "success",
			Message: "Register successful",
			Data:    userInfo,
		})
    }
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token, err := c.Cookie("access_token")
        if err != nil || token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }

        userID, err := usecase.ValidateToken(token) 
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        c.Set("userID", userID) 
        c.Next()
    }
}


func GetMe(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }

        user, err := userRepo.FindByID(userID.(int))
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
            return
        }

        c.JSON(http.StatusOK, APIResponse{
            Status:  "success",
            Message: "Success",
            Data:    user,
        })

    }
}



func RegisterAuthRoutes(router *gin.Engine, userRepo repository.UserRepository, ch *amqp091.Channel) {
    api := router.Group("/api/v1/auth")
    {
        api.POST("/login", Login(userRepo))
        api.POST("/register", Register(userRepo, ch))
        api.GET("/me",AuthMiddleware(), GetMe(userRepo))
    }
}