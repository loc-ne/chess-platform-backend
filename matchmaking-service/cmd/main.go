package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/locne/matchmaking-service/internal/interface/handler"
    "github.com/locne/matchmaking-service/internal/usecase"
    "github.com/locne/matchmaking-service/internal/infrastructure/messagebroker"
    "github.com/gin-contrib/cors"
    "log"
    "fmt"
)

func main() {
    if os.Getenv("ENV") != "production" {
        godotenv.Load("internal/infrastructure/config/.env")
    }
    poolManager := usecase.NewPoolManager()

    mqConn, mqCh, err := messagebroker.ConnectRabbit()
    if err != nil {
        log.Fatalf("RabbitMQ connection error: %v", err)
    }
    defer mqConn.Close()
    defer mqCh.Close()

    workerPool := usecase.NewWorkerPool(poolManager, 3, mqCh)

    router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    handler.RegisterMatchmakingRoutes(router, workerPool, poolManager)
    port := os.Getenv("PORT")
    fmt.Println("Server started on port:", port)
    if runErr := router.Run(":" + port); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}