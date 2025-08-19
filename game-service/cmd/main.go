package main

import (
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/go-redis/redis/v8"
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "github.com/locne/game-service/internal/usecase/game"
    "github.com/locne/game-service/internal/infrastructure/messagebroker"
    "github.com/locne/game-service/internal/infrastructure/db"
    "github.com/locne/game-service/internal/interface/repository"
)

func main() {
    // Load environment variables
    if os.Getenv("ENV") != "production" {
        godotenv.Load("internal/infrastructure/config/.env")
    }
    
    // Setup Gin router
    frontendEnv := os.Getenv("FRONTEND_ORIGIN")
    router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", frontendEnv},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Setup Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "redis-15123.crce185.ap-seast-1-1.ec2.redns.redis-cloud.com:15123",
        Username: "default",
        Password: "djPg9bQavZmVK6SH2e5npv5DBwolxLP6",
        DB:       0,
    })

    // Setup MongoDB connection
    mongoDB, err := db.ConnectMongoDB()
    if err != nil {
        log.Fatalf("MongoDB connection error: %v", err)
    }
    defer mongoDB.Disconnect()

    // Setup repository
    gameRepo := repository.NewGameRepository(mongoDB.Database)

    // Setup context
    ctx := context.Background()

    // Setup GameManager with MongoDB repository
    gameManager := game.NewGameManager(redisClient, ctx, gameRepo)

    // Setup RabbitMQ
    mqConn, mqCh, err := messagebroker.ConnectRabbit()
    if err != nil {
        log.Fatalf("RabbitMQ connection error: %v", err)
    }
    defer mqConn.Close()
    defer mqCh.Close()

    // Start consuming game creation messages
    messagebroker.ConsumeGameCreate(mqCh, gameManager)

    // Start listening for moves and game actions
    go gameManager.ListenChannels()

    // Setup graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh

        log.Println(" Shutting down gracefully...")
        
        // Close connections
        if err := mongoDB.Disconnect(); err != nil {
            log.Printf("Error disconnecting MongoDB: %v", err)
        }
        
        if err := redisClient.Close(); err != nil {
            log.Printf("Error closing Redis: %v", err)
        }
        
        os.Exit(0)
    }()
    
    port := os.Getenv("PORT")
    fmt.Println("Server started on port:", port)
    if runErr := router.Run(":" + port); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}