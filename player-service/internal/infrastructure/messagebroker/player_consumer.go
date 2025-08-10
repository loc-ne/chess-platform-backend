package messagebroker

import (
    "encoding/json"
    "log"

    "github.com/locne/player-service/internal/usecase"
    "github.com/locne/player-service/internal/interface/repository"
    "github.com/rabbitmq/amqp091-go"
)

type RegisterMsg struct {
    UserID int `json:"user_id"`
    Elo    int `json:"elo"`
}

func ConsumePlayerRegister(ch *amqp091.Channel, repo repository.PlayerRepository) {
    msgs, err := ch.Consume(
        "user.register", // queue
        "",              // consumer
        true,            // auto-ack
        false,           // exclusive
        false,           // no-local
        false,           // no-wait
        nil,             // args
    )
    if err != nil {
        log.Fatalf("Failed to register consumer: %v", err)
    }

    go func() {
        for d := range msgs {
            var msg RegisterMsg
            if err := json.Unmarshal(d.Body, &msg); err != nil {
                log.Printf("Invalid message: %v", err)
                continue
            }
            if err := usecase.CreateAllGameTypePlayers(repo, msg.UserID, msg.Elo); err != nil {
                log.Printf("Create player error: %v", err)
            }
        }
    }()
}