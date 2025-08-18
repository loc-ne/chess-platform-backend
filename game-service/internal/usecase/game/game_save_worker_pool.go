package game

import (
    "context"
    "log"
    "time"
    "github.com/locne/game-service/internal/entity"
    "github.com/locne/game-service/internal/interface/repository"
)

type GameSaveJob struct {
    Game      entity.Game
    Timestamp time.Time
}

type GameSaveWorkerPool struct {
    jobs       chan GameSaveJob
    batchSize  int
    flushTime  time.Duration
    repo       repository.GameRepository
    workers    int
}

func NewGameSaveWorkerPool(repo repository.GameRepository, workers int) *GameSaveWorkerPool {
    pool := &GameSaveWorkerPool{
        jobs:      make(chan GameSaveJob, 1000), // Buffer for 1000 games
        batchSize: 50,                           // Batch size
        flushTime: 5 * time.Second,             // Max wait time
        repo:      repo,
        workers:   workers,
    }
    
    // Start workers
    for i := 0; i < workers; i++ {
        go pool.worker()
    }
    
    return pool
}

func (pool *GameSaveWorkerPool) worker() {
    ticker := time.NewTicker(pool.flushTime)
    defer ticker.Stop()
    
    batch := make([]entity.Game, 0, pool.batchSize)
    
    for {
        select {
        case job := <-pool.jobs:
            batch = append(batch, job.Game)
            
            if len(batch) >= pool.batchSize {
                pool.flushBatch(batch)
                batch = batch[:0] 
            }
            
        case <-ticker.C:
            if len(batch) > 0 {
                pool.flushBatch(batch)
                batch = batch[:0]
            }
        }
    }
}

func (pool *GameSaveWorkerPool) flushBatch(games []entity.Game) {
    if len(games) == 0 {
        return
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    start := time.Now()
    err := pool.repo.SaveGamesBatch(ctx, games)
    duration := time.Since(start)
    
    if err != nil {
        log.Printf("Failed to save batch of %d games: %v (took %v)", 
            len(games), err, duration)
        
        for _, game := range games {
            log.Printf("Failed game: %s", game.GameID)
        }
    } else {
        log.Printf("Successfully saved batch of %d games (took %v)", 
            len(games), duration)
    }
}

func (pool *GameSaveWorkerPool) SaveGame(game entity.Game) {
    select {
    case pool.jobs <- GameSaveJob{Game: game, Timestamp: time.Now()}:
    default:
        log.Printf("Game save queue full, dropping game %s", game.GameID)
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            if err := pool.repo.SaveGame(ctx, game); err != nil {
                log.Printf("Failed to save game %s: %v", game.GameID, err)
            }
        }()
    }
}