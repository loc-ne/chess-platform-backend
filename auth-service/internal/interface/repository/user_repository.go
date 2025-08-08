package repository

import (
    "errors"
    "github.com/locne/auth-service/internal/entity"
)

type UserRepository interface {
    FindByEmail(email string) (*User, error)
    FindByUsername(username string) (*User, error)
    Create(user *User) error
    SaveRefreshToken(userID uint, refreshToken string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
    var user User
    result := r.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
    var user User
    result := r.db.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (r *userRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) SaveRefreshToken(userID uint, refreshToken string) error {
    return r.db.Model(&User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}