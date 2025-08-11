package usecase

import (
    "github.com/locne/player-service/internal/entity"
    "github.com/locne/player-service/internal/interface/repository"
)

func CreateAllGameTypePlayers(repo repository.PlayerRepository, userID int, elo int) error {
    gameTypes := []entity.GameType{
        entity.GameTypeBullet,
        entity.GameTypeBlitz,
        entity.GameTypeRapid,
        entity.GameTypeClassical,
    }

    for _, gt := range gameTypes {
        player := entity.Player{
            UserID:      userID,
            GameType:    gt,
            WhiteGames:  0,
            BlackGames:  0,
            Rating:      elo,
            GamesPlayed: 0,
            Wins:        0,
            Losses:      0,
            Draws:       0,
            PeakRating:  elo,
        }
        if err := repo.Create(player); err != nil {
            return err
        }
    }
    return nil
}

func GetPlayerEloByGameType(repo repository.PlayerRepository, userID int, gameType entity.GameType) (int, error) {
    player, err := repo.GetByUserIDAndGameType(userID, gameType)
    if err != nil {
        return 0, err
    }
    return player.Rating, nil
}