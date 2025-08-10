package messagebroker

import (
    "encoding/json"
    "fmt"
    "github.com/rabbitmq/amqp091-go"
)

type PlayerRegisterMsg struct {
    UserID int `json:"user_id"`
    Elo    int `json:"elo"`
}

func PublishPlayerRegister(ch *amqp091.Channel, userID int, elo int) error {
    msg := PlayerRegisterMsg{
        UserID: userID,
        Elo:    elo,
    }
    body, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("marshal message error: %v", err)
    }

    err = ch.Publish(
        "",              // exchange
        "user.register", // routing key (queue name)
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
