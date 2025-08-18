package messagebroker

import (
    "encoding/json"
    "log"
    "time"
    "github.com/locne/game-service/internal/usecase/game"
    "github.com/locne/game-service/internal/usecase/engine"
    "github.com/rabbitmq/amqp091-go"
    "fmt"
    "math/rand"
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

type TimeControl struct {
    Type        string `json:"type"`
    InitialTime int    `json:"initialTime"`
    Increment   int    `json:"increment"`
}

type CreateGameMsg struct {
    Player1     PlayerGameInfo `json:"player1"`
    Player2     PlayerGameInfo `json:"player2"`
    TimeControl TimeControl    `json:"timeControl"`
    Colors      Colors         `json:"colors"`
}

func ConsumeGameCreate(ch *amqp091.Channel, gm *game.GameManager) {
    msgs, err := ch.Consume(
        "game.create", // queue
        "",            // consumer
        true,          // auto-ack
        false,         // exclusive
        false,         // no-local
        false,         // no-wait
        nil,           // args
    )
    if err != nil {
        log.Fatalf("Failed to register consumer: %v", err)
    }

    go func() {
        for d := range msgs {
            var msg CreateGameMsg
            if err := json.Unmarshal(d.Body, &msg); err != nil {
                log.Printf("Invalid message: %v", err)
                continue
            }
            
            if err := createGameFromMessage(msg, gm); err != nil {
                log.Printf("Create game error: %v", err)
            }
        }
    }()
}

func createGameFromMessage(msg CreateGameMsg, gm *game.GameManager) error {
    gameID := generateGameID()
    
    players := make(map[int]*game.Player)
    
    players[msg.Player1.UserID] = &game.Player{
        ID:       msg.Player1.UserID,
        Username: msg.Player1.Username,
        Color:    msg.Colors.Player1,
        Rating:   msg.Player1.Rating,
        IsOnline: true,
    }
    
    players[msg.Player2.UserID] = &game.Player{
        ID:       msg.Player2.UserID,
        Username: msg.Player2.Username,
        Color:    msg.Colors.Player2,
        Rating:   msg.Player2.Rating,
        IsOnline: true,
    }
    
    chessEngine := &engine.ChessEngine{}
    initialGameState := chessEngine.CreateServerGameState()
    
    timeControl := &game.TimeControl{
        Type:        msg.TimeControl.Type,
        InitialTime: msg.TimeControl.InitialTime,
        Increment:   msg.TimeControl.Increment,
    }
    
    newGame := &game.Game{
        ID:            gameID,
        Players:       players,
        Spectators:    make(map[int]*game.Player),
        GameState:     initialGameState,
        TimeControl:   timeControl,
        WhiteTimeLeft: msg.TimeControl.InitialTime,
        BlackTimeLeft: msg.TimeControl.InitialTime,
        LastMoveTime:  time.Now(),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
        DrawOffers:    make(map[string]*game.DrawOffer), // Initialize draw offers map
    }
    
    gm.AddGame(newGame)
    
    log.Printf("✅ Created game %s: %s (%s) vs %s (%s)", 
        gameID,
        msg.Player1.Username, msg.Colors.Player1,
        msg.Player2.Username, msg.Colors.Player2,
    )
    
    // Chỉ publish matchFound để client navigate, không gửi game state
    matchFoundMsg := game.StateUpdateMessage{
        Type:    "matchFound",
        RoomID:  gameID,
        Player1: *players[msg.Player1.UserID],
        Player2: *players[msg.Player2.UserID],
    }
    
    gm.PublishStateUpdate(matchFoundMsg)
    
    return nil
}

func generateGameID() string {
    return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(10000))
}