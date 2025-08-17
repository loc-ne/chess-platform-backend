package usecase

import (
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
    "time"
    "context"
    "github.com/locne/matchmaking-service/internal/infrastructure/messagebroker"
    "github.com/rabbitmq/amqp091-go"
    "math/rand"
)

var httpClient = &http.Client{
    Timeout: 10 * time.Second, // Overall request timeout
    Transport: &http.Transport{
        MaxIdleConns:          100,              // Maximum idle connections across all hosts
        MaxIdleConnsPerHost:   20,               // Maximum idle connections per host
        IdleConnTimeout:       90 * time.Second, // How long idle connections stay alive
        TLSHandshakeTimeout:   10 * time.Second, // TLS handshake timeout
        ExpectContinueTimeout: 1 * time.Second,  // Expect: 100-continue timeout
        DisableKeepAlives:     false,            // Enable keep-alive
        DisableCompression:    false,            // Enable compression
        MaxConnsPerHost:       50,               // Maximum connections per host
        ResponseHeaderTimeout: 10 * time.Second, // Response header timeout
    },
}

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

func makeHTTPRequest(url string, payload map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
    bodyBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("marshal error: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        return nil, fmt.Errorf("create request error: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Connection", "keep-alive") 

    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request error: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
    }

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response error: %w", err)
    }

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        return nil, fmt.Errorf("unmarshal error: %w", err)
    }

    return result, nil
}

func getPlayerElo(playerServiceURL string, userID int, gameType string) (int, error) {
    payload := map[string]interface{}{
        "user_id":   userID,
        "game_type": gameType,
    }

    result, err := makeHTTPRequest(
        playerServiceURL+"/api/v1/player/elo", 
        payload, 
        5*time.Second, 
    )
    if err != nil {
        return 0, err
    }

    elo, ok := result["elo"].(float64)
    if !ok {
        return 0, fmt.Errorf("elo not found in response")
    }

    return int(elo), nil
}

func getPlayerColorBalance(playerServiceURL string, userID int, gameType string) (float64, error) {
    payload := map[string]interface{}{
        "user_id":   userID,
        "game_type": gameType,
    }

    result, err := makeHTTPRequest(
        playerServiceURL+"/api/v1/player/color_balance", 
        payload, 
        3*time.Second, 
    )
    if err != nil {
        return 0, err
    }

    balance, ok := result["colorBalance"].(float64)
    if !ok {
        return 0, fmt.Errorf("colorBalance not found in response")
    }

    return balance, nil
}

func AssignUserColors(player1ID int, player2ID int, gameType string, playerServiceURL string) (player1Color, player2Color string, err error) {
    type colorResult struct {
        balance float64
        err     error
    }

    p1Result := make(chan colorResult, 1)
    p2Result := make(chan colorResult, 1)

    go func() {
        balance, err := getPlayerColorBalance(playerServiceURL, player1ID, gameType)
        p1Result <- colorResult{balance: balance, err: err}
    }()

    go func() {
        balance, err := getPlayerColorBalance(playerServiceURL, player2ID, gameType)
        p2Result <- colorResult{balance: balance, err: err}
    }()

    ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
    defer cancel()

    var p1Balance, p2Balance float64

    select {
    case r1 := <-p1Result:
        if r1.err != nil {
            return "", "", fmt.Errorf("player1 color balance error: %w", r1.err)
        }
        p1Balance = r1.balance
    case <-ctx.Done():
        return "", "", fmt.Errorf("timeout getting player1 color balance")
    }

    select {
    case r2 := <-p2Result:
        if r2.err != nil {
            return "", "", fmt.Errorf("player2 color balance error: %w", r2.err)
        }
        p2Balance = r2.balance
    case <-ctx.Done():
        return "", "", fmt.Errorf("timeout getting player2 color balance")
    }

    // Assign colors based on balance
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
        
        // Get ELO với fallback
        elo, err := getPlayerElo(playerServiceURL, job.Player.UserId, job.TimeControl.Type)
        if err != nil {
            fmt.Printf("Failed to get ELO for user %d: %v, using default 1200\n", job.Player.UserId, err)
            job.Player.Elo = 1200 
        } else {
            job.Player.Elo = elo
        }

        opponent := wp.PoolManager.FindNearestElo(job.PoolKey, job.Player.Elo)
        if opponent != nil {
            p1Color, p2Color, err := AssignUserColors(job.Player.UserId, opponent.UserId, job.TimeControl.Type, playerServiceURL)
            if err != nil {
                fmt.Printf("Assign color error: %v, using random assignment\n", err)
                if rand.Float64() < 0.5 {
                    p1Color, p2Color = "white", "black"
                } else {
                    p1Color, p2Color = "black", "white"
                }
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
            } else {
                wp.PoolManager.Leave(job.PoolKey, opponent.UserId)
                wp.PoolManager.Leave(job.PoolKey, job.Player.UserId)
                fmt.Printf("Match created: Player %d vs Player %d, both players removed from pool\n", 
                    job.Player.UserId, opponent.UserId)
            }
        } else {
            wp.PoolManager.Join(job.PoolKey, job.Player)
        }
    }
}