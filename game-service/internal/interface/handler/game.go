package handler

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/locne/game-service/internal/usecase"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type GameHandler struct {
    gameUseCase *usecase.GameUseCase
}

func NewGameHandler(gameUseCase *usecase.GameUseCase) *GameHandler {
    return &GameHandler{
        gameUseCase: gameUseCase,
    }
}

func (h *GameHandler) GetGameByID(c *gin.Context) {
    gameID := c.Param("gameId")
    
    if gameID == "" {
        c.JSON(http.StatusBadRequest, APIResponse{
            Status:  "error",
            Message: "Game ID is required",
            Error:   "missing gameId parameter",
        })
        return
    }
    
    game, err := h.gameUseCase.GetGameByID(c.Request.Context(), gameID)
    if err != nil {
        if strings.Contains(err.Error(), "game not found") {
            c.JSON(http.StatusNotFound, APIResponse{
                Status:  "error",
                Message: "Game not found",
                Error:   err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, APIResponse{
            Status:  "error",
            Message: "Failed to retrieve game",
            Error:   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, APIResponse{
        Status:  "success",
        Message: "Game retrieved successfully",
        Data:    game,
    })
}

func RegisterGameRoutes(router *gin.Engine, gameUseCase *usecase.GameUseCase) {
    gameHandler := NewGameHandler(gameUseCase)
    
    api := router.Group("/api/v1")
    {
        api.GET("/games/:gameId", gameHandler.GetGameByID)
    }
}