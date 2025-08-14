package main

import (
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/go-redis/redis/v8"
    "context"
    "log"
    "github.com/locne/game-service/internal/usecase/game"
    "github.com/locne/game-service/internal/infrastructure/messagebroker"
)

func main() {
    godotenv.Load("internal/infrastructure/config/.env")
    router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    redisClient := redis.NewClient(&redis.Options{
        Addr:     "redis-15123.crce185.ap-seast-1-1.ec2.redns.redis-cloud.com:15123",
        Username: "default",
        Password: "djPg9bQavZmVK6SH2e5npv5DBwolxLP6",
        DB:       0,
    })

    ctx := context.Background()

    gameManager := game.NewGameManager(redisClient, ctx)

    mqConn, mqCh, err := messagebroker.ConnectRabbit()
    if err != nil {
        log.Fatalf("RabbitMQ connection error: %v", err)
    }
    defer mqConn.Close()
    defer mqCh.Close()

    messagebroker.ConsumeGameCreate(mqCh, gameManager)

    go gameManager.ListenMoves()

    router.Run(":3004")
}