package handler

import (
    "net/http"
    "github.com/locne/matchmaking-service/internal/usecase"
    "github.com/gin-gonic/gin"
    "fmt"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  []string    `json:"errors,omitempty"`
}

type FindMatchDto struct {
    TimeControl usecase.TimeControl `json:"timeControl"`
    Player      usecase.Player      `json:"player"`
}


func JoinMatchmakingPool(workerPool *usecase.WorkerPool) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req FindMatchDto
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        key := fmt.Sprintf("%d_%d", req.TimeControl.InitialTime, req.TimeControl.Increment)
        workerPool.Jobs <- usecase.MatchmakingJob{
            PoolKey: key,
            Player:  req.Player,
            TimeControl: req.TimeControl,
        }
        c.JSON(http.StatusOK, APIResponse{
            Status:  "waiting",
            Message: "Added to pool, waiting for opponent",
            Data:    req.Player,
        })
        
    }
}

func LeaveMatchmakingPool(poolManager *usecase.PoolManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req FindMatchDto
        if err := c.ShouldBindJSON(&req); err != nil {
            fmt.Println("BindJSON error:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        key := fmt.Sprintf("%d_%d", req.TimeControl.InitialTime, req.TimeControl.Increment)
        poolManager.Leave(key, req.Player.UserId)
        c.JSON(http.StatusOK, APIResponse{
            Status:  "success",
            Message: "Player removed from matchmaking pool",
            Data:    req.Player,
        })
    }
}

func RegisterMatchmakingRoutes(router *gin.Engine, workerPool *usecase.WorkerPool, poolManager *usecase.PoolManager) {
    api := router.Group("/api/v1/matchmaking")
    {
        api.POST("/join", JoinMatchmakingPool(workerPool))
        api.POST("/leave", LeaveMatchmakingPool(poolManager))
    }
}