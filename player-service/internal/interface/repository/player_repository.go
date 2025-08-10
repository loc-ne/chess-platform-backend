package repository

import (
    "github.com/locne/player-service/internal/entity"
    "gorm.io/gorm"
)

type PlayerRepository interface {
    Create(player entity.Player) error
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