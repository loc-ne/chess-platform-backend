package usecase

import (
    "fmt"
)

type MatchmakingJob struct {
    PoolKey string
    Player  Player
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

func (wp *WorkerPool) worker() {
    for job := range wp.Jobs {
        opponent := wp.PoolManager.FindNearestElo(job.PoolKey, job.Player.Elo)
        if opponent != nil {
            // to do MQ
        } else {
            wp.PoolManager.Join(job.PoolKey, job.Player)
            job.Done <- nil
        }
    }
}