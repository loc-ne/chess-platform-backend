package repository

import (
    "github.com/locne/player-service/internal/entity"
    "gorm.io/gorm"
)

type PlayerRepository interface {
    Create(player entity.Player) error
    GetByUserIDAndGameType(userID int, gameType entity.GameType) (entity.Player, error)
}

type playerRepository struct {
    db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
    return &playerRepository{db: db}
}

func (r *playerRepository) Create(player entity.Player) error {
    return r.db.Create(&player).Error
}

func (r *playerRepository) GetByUserIDAndGameType(userID int, gameType entity.GameType) (entity.Player, error) {
    var player entity.Player
    err := r.db.Where("user_id = ? AND game_type = ?", userID, gameType).First(&player).Error
    return player, err
}