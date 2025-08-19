package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/locne/player-service/internal/infrastructure/db"
    "github.com/locne/player-service/internal/interface/repository"
    "github.com/locne/player-service/internal/interface/handler"
    "github.com/locne/player-service/internal/entity"
	"github.com/locne/player-service/internal/infrastructure/messagebroker"
    "github.com/gin-contrib/cors"
    "os"
)

func main() {
    if os.Getenv("ENV") != "production" {
        godotenv.Load("internal/infrastructure/config/.env")
    }
    router := gin.Default()

    frontendEnv := os.Getenv("FRONTEND_ORIGIN")
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", frontendEnv},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    dbConn, err := db.ConnectPostgres()
    if err != nil {
        panic(err)
    }

    dbConn.AutoMigrate(&entity.Player{})
    playerRepository := repository.NewPlayerRepository(dbConn)

	conn, ch, err := messagebroker.ConnectRabbit()
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer ch.Close()
    messagebroker.ConsumePlayerRegister(ch, playerRepository)

    handler.RegisterPlayerRoutes(router, playerRepository)

    port := os.Getenv("PORT")
    fmt.Println("Server started on port:", port)
    if runErr := router.Run(":" + port); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}