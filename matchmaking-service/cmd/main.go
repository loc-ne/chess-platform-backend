package main

import (
    "github.com/gin-gonic/gin"
    "github.com/locne/matchmaking-service/internal/interface/handler"
    "github.com/locne/matchmaking-service/internal/usecase"
	"github.com/gin-contrib/cors"
)

func main() {
    poolManager := usecase.NewPoolManager()

    workerPool := usecase.NewWorkerPool(poolManager, 3)

    router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    handler.RegisterMatchmakingRoutes(router, workerPool, poolManager)

    router.Run(":3003")
}