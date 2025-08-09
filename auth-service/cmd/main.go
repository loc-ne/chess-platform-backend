package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/locne/auth-service/internal/infrastructure/db"
    "github.com/locne/auth-service/internal/interface/repository"
    "github.com/locne/auth-service/internal/interface/handler"
    "github.com/locne/auth-service/internal/entity"
    "github.com/gin-contrib/cors"
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

    dbConn, err := db.ConnectPostgres()
    if err != nil {
        panic(err)
    }

    dbConn.AutoMigrate(&entity.User{})
    userRepository := repository.NewUserRepository(dbConn)
    handler.RegisterAuthRoutes(router, userRepository)

    fmt.Println("Server started on :3001")
    if runErr := router.Run(":3001"); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}