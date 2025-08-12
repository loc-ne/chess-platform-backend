package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/locne/player-service/internal/interface/repository"
    "github.com/locne/player-service/internal/entity"
)

type GetEloRequest struct {
    UserID   int                `json:"user_id"`
    GameType entity.GameType    `json:"game_type"`
}

type GetColorBalanceRequest struct {
    UserID   int             `json:"user_id"`
    GameType entity.GameType `json:"game_type"`
}

func GetPlayerElo(repo repository.PlayerRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req GetEloRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        player, err := repo.GetByUserIDAndGameType(req.UserID, req.GameType)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "user_id": req.UserID,
            "game_type": req.GameType,
            "elo": player.Rating,
        })
    }
}

func GetPlayerColorBalance(repo repository.PlayerRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req GetColorBalanceRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        player, err := repo.GetByUserIDAndGameType(req.UserID, req.GameType)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "colorBalance": player.ColorBalance(),
        })
    }
}

func RegisterPlayerRoutes(router *gin.Engine, repo repository.PlayerRepository) {
    api := router.Group("/api/v1/player")
    {
        api.POST("/elo", GetPlayerElo(repo))
        api.POST("/color_balance", GetPlayerColorBalance(repo))
    }
}