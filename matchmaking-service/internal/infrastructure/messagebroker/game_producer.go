package messagebroker

import (
    "encoding/json"
    "fmt"
    "github.com/rabbitmq/amqp091-go"
    "github.com/locne/matchmaking-service/internal/usecase"
)

type PlayerGameInfo struct {
    UserID   int    `json:"userId"`
    Username string `json:"username"`
    Rating   int    `json:"rating"`
}

type Colors struct {
    Player1 string `json:"player1"` 
    Player2 string `json:"player2"` 
}


type CreateGameMsg struct {
    Player1     PlayerGameInfo   `json:"player1"`
    Player2     PlayerGameInfo   `json:"player2"`
    TimeControl usecase.TimeControl `json:"timeControl"`
    Colors      Colors          `json:"colors"`
}

func PublishGameCreate(ch *amqp091.Channel, msg CreateGameMsg) error {
    body, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("marshal message error: %v", err)
    }

    err = ch.Publish(
        "",              // exchange
        "game.create",   // routing key (queue name)
        false,           // mandatory
        false,           // immediate
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        return fmt.Errorf("publish message error: %v", err)
    }
    return nil
}
