package entity

import "time"

type User struct {
    ID           int       `gorm:"primaryKey"`
    Email        string    `gorm:"uniqueIndex;size:255"`
    Username     string    `gorm:"uniqueIndex;size:100"`
    Password     string    `gorm:"size:255"`
    IsActive     bool      `gorm:"default:true"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    RefreshToken string    `gorm:"size:255"`
}