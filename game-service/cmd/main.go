package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
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

    
}