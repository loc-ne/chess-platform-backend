package usecase

import (
    "fmt"
    "net/http"
    "encoding/json"
    "bytes"
    "io/ioutil"
)

type MatchmakingJob struct {
    PoolKey string
    Player  Player
    TimeControl TimeControl
}

type WorkerPool struct {
    Jobs chan MatchmakingJob
    PoolManager *PoolManager
}

func NewWorkerPool(poolManager *PoolManager, workerCount int) *WorkerPool {
    wp := &WorkerPool{
        Jobs: make(chan MatchmakingJob, 1000),
        PoolManager: poolManager,
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

func (wp *WorkerPool) worker() {
    for job := range wp.Jobs {
        playerServiceURL := "http://localhost:3002"
        elo, err := getPlayerElo(playerServiceURL, job.Player.UserId, job.TimeControl.Type)
        if err == nil {
            job.Player.Elo = elo
        }
        opponent := wp.PoolManager.FindNearestElo(job.PoolKey, job.Player.Elo)
        if opponent != nil {
            
        } else {
            wp.PoolManager.Join(job.PoolKey, job.Player)
        }
    }
}