package messagebroker

import (
    "fmt"
    "os"
    "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbit() (*amqp091.Connection, *amqp091.Channel, error) {
    url := os.Getenv("RABBITMQ_URL")
    if url == "" {
        return nil, nil, fmt.Errorf("RABBITMQ_URL env not set")
    }

    conn, err := amqp091.Dial(url)
    if err != nil {
        return nil, nil, fmt.Errorf("Can't connect to RabbitMQ: %v", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, nil, fmt.Errorf("Can't open channel: %v", err)
    }

    _, err = ch.QueueDeclare(
        "game.create", // name
        true,         // durable
        false,        // autoDelete
        false,        // exclusive
        false,        // noWait
        nil,          // arguments
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return nil, nil, fmt.Errorf("Can't declare queue: %v", err)
    }

    fmt.Println("RabbitMQ connected, queue declared: user.register")
    return conn, ch, nil
}
