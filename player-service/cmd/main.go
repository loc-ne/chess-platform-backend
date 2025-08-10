package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "github.com/locne/player-service/internal/infrastructure/db"
    "github.com/locne/player-service/internal/interface/repository"
    "github.com/locne/player-service/internal/entity"
	"github.com/locne/player-service/internal/infrastructure/messagebroker"
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

    dbConn.AutoMigrate(&entity.Player{})
    playerRepository := repository.NewPlayerRepository(dbConn)

	conn, ch, err := messagebroker.ConnectRabbit()
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    defer ch.Close()
    messagebroker.ConsumePlayerRegister(ch, playerRepository)

    fmt.Println("Server started on :3002")
    if runErr := router.Run(":3002"); runErr != nil {
        panic("ListenAndServe: " + runErr.Error())
    }
}