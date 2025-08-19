package main

import (
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/go-redis/redis/v8"
    "github.com/locne/ws-service/internal/interface/handler"
    "github.com/locne/ws-service/internal/usecase"
    "fmt"
    "os"
)

func main() {
    if os.Getenv("ENV") != "production" {
        godotenv.Load("internal/infrastructure/config/.env")
    }

    redisClient := redis.NewClient(&redis.Options{
        Addr:     "redis-15123.crce185.ap-seast-1-1.ec2.redns.redis-cloud.com:15123",
        Username: "default",
        Password: "djPg9bQavZmVK6SH2e5npv5DBwolxLP6",
        DB:       0,
    })
    

    roomManager := usecase.NewRoomManager(redisClient)
    go roomManager.ListenStateUpdates()
    
    router := gin.Default()
    
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
    
    handler.RegisterWebSocketRoutes(router, roomManager)
    
    port := os.Getenv("PORT")
    fmt.Println("Server started on port:", port)
    if runErr := router.Run(":" + port); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}