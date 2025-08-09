package entity

import (
    "time"
)

type GameType string

const (
    GameTypeBullet    GameType = "bullet"
    GameTypeBlitz     GameType = "blitz"
    GameTypeRapid     GameType = "rapid"
    GameTypeClassical GameType = "classical"
)

type Player struct {
    ID          int       `gorm:"primaryKey"`
    UserID      int       `gorm:"index"`
    GameType    GameType  `gorm:"type:enum('bullet','blitz','rapid','classical');index"`
    WhiteGames  int       `gorm:"default:0"`
    BlackGames  int       `gorm:"default:0"`
    Rating      int       `gorm:"default:1200"`
    GamesPlayed int       `gorm:"default:0"`
    Wins        int       `gorm:"default:0"`
    Losses      int       `gorm:"default:0"`
    Draws       int       `gorm:"default:0"`
    PeakRating  int       `gorm:"default:1200"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (p *Player) WinRate() int {
    if p.GamesPlayed > 0 {
        return int(float64(p.Wins) / float64(p.GamesPlayed) * 100)
    }
    return 0
}

func (p *Player) ColorBalance() float64 {
    totalGames := p.WhiteGames + p.BlackGames
    if totalGames == 0 {
        return 0
    }
    return float64(p.WhiteGames-p.BlackGames) / float64(totalGames)
}