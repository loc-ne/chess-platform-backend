package usecase

import (
    "context"
    "fmt"
    
    "github.com/locne/game-service/internal/entity"
    "github.com/locne/game-service/internal/interface/repository"
)

type GameUseCase struct {
    gameRepo repository.GameRepository
}

func NewGameUseCase(gameRepo repository.GameRepository) *GameUseCase {
    return &GameUseCase{
        gameRepo: gameRepo,
    }
}

func (uc *GameUseCase) GetGameByID(ctx context.Context, gameID string) (*entity.Game, error) {
    if gameID == "" {
        return nil, fmt.Errorf("gameID cannot be empty")
    }
    
    game, err := uc.gameRepo.GetGameByID(ctx, gameID)
    if err != nil {
        return nil, err
    }
    
    return game, nil
}