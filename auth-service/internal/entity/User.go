package entity

import "time"

type User struct {
	ID           int
	Email        string
	Username     string
	Password     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RefreshToken string
}
