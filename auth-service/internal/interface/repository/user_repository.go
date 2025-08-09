package repository

import (
    "github.com/locne/auth-service/internal/entity"
    "gorm.io/gorm"
)

type UserRepository interface {
    FindByEmail(email string) (entity.User, error)
    FindByUsername(username string) (entity.User, error)
    Create(user entity.User) error
    SaveRefreshToken(userID int, refreshToken string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
    var user entity.User
    result := r.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return user, result.Error
    }
    return user, nil
}

func (r *userRepository) FindByUsername(username string) (entity.User, error) {
    var user entity.User
    result := r.db.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return user, result.Error
    }
    return user, nil
}

func (r *userRepository) Create(user entity.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) SaveRefreshToken(userID int, refreshToken string) error {
    return r.db.Model(entity.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}