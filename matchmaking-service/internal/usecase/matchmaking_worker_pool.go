package usecase

import (
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
    "github.com/locne/matchmaking-service/internal/infrastructure/messagebroker"
    "github.com/rabbitmq/amqp091-go"
    "math/rand"
)

type MatchmakingJob struct {
    PoolKey string
    Player  Player
    TimeControl TimeControl
}

type WorkerPool struct {
    Jobs chan MatchmakingJob
    PoolManager *PoolManager
    MQChannel   *amqp091.Channel
}

func NewWorkerPool(poolManager *PoolManager, workerCount int, mqCh *amqp091.Channel) *WorkerPool {
    wp := &WorkerPool{
        Jobs: make(chan MatchmakingJob, 1000),
        PoolManager: poolManager,
        MQChannel:   mqCh,
    }
    for i := 0; i < workerCount; i++ {
        go wp.worker()
    }
    return wp
}

func getPlayerElo(playerServiceURL string, userID int, gameType string) (int, error) {
    reqBody := map[string]interface{}{
        "user_id":  userID,
        "game_type": gameType,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(playerServiceURL+"/api/v1/player/elo", "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("player-service returned status %d", resp.StatusCode)
    }
    respBody, _ := ioutil.ReadAll(resp.Body)
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return 0, err
    }
    elo, ok := result["elo"].(float64)
    if !ok {
        return 0, fmt.Errorf("elo not found in response")
    }
    return int(elo), nil
}

func getPlayerColorBalance(playerServiceURL string, userID int, gameType string) (float64, error) {
    reqBody := map[string]interface{}{
        "user_id":  userID,
        "game_type": gameType,
    }
    bodyBytes, _ := json.Marshal(reqBody)
    resp, err := http.Post(playerServiceURL+"/api/v1/player/color_balance", "application/json", bytes.NewBuffer(bodyBytes))
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("player-service returned status %d", resp.StatusCode)
    }
    respBody, _ := ioutil.ReadAll(resp.Body)
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return 0, err
    }
    balance, ok := result["colorBalance"].(float64)
    if !ok {
        return 0, fmt.Errorf("colorBalance not found in response")
    }
    return balance, nil
}

func AssignUserColors(player1ID int, player2ID int, gameType string, playerServiceURL string) (player1Color, player2Color string, err error) {
    p1Balance, err := getPlayerColorBalance(playerServiceURL, player1ID, gameType)
    if err != nil {
        return "", "", err
    }
    p2Balance, err := getPlayerColorBalance(playerServiceURL, player2ID, gameType)
    if err != nil {
        return "", "", err
    }

    if p1Balance < p2Balance {
        return "white", "black", nil
    } else if p2Balance < p1Balance {
        return "black", "white", nil
    } else {
        if rand.Float64() < 0.5 {
            return "white", "black", nil
        }
        return "black", "white", nil
    }
}

func (wp *WorkerPool) worker() {
    for job := range wp.Jobs {
        playerServiceURL := "http://localhost:3002"
        elo, err := getPlayerElo(playerServiceURL, job.Player.UserId, job.TimeControl.Type)
        if err == nil {
            job.Player.Elo = elo
        }
        opponent := wp.PoolManager.FindNearestElo(job.PoolKey, job.Player.Elo)
        if opponent != nil {
        p1Color, p2Color, err := AssignUserColors(job.Player.UserId, opponent.UserId, job.TimeControl.Type, playerServiceURL)
        if err != nil {
            fmt.Println("Assign color error:", err)
            continue
        }
        gameMsg := messagebroker.CreateGameMsg{
            Player1: messagebroker.PlayerGameInfo{
                UserID:   job.Player.UserId,
                Username: job.Player.UserName,
                Rating:   job.Player.Elo,
            },
            Player2: messagebroker.PlayerGameInfo{
                UserID:   opponent.UserId,
                Username: opponent.UserName,
                Rating:   opponent.Elo,
            },
            TimeControl: messagebroker.TimeControl{
            Type:        job.TimeControl.Type,
            InitialTime: job.TimeControl.InitialTime,
            Increment:   job.TimeControl.Increment,
            },
            Colors: messagebroker.Colors{
                Player1: p1Color,
                Player2: p2Color,
            },
        }
        err = messagebroker.PublishGameCreate(wp.MQChannel, gameMsg)
        if err != nil {
            fmt.Println("publish game.create error:", err)
        }
        } else {
                wp.PoolManager.Join(job.PoolKey, job.Player)
            }
        }
}