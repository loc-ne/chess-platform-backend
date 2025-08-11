package usecase

import (
    "github.com/locne/auth-service/internal/interface/repository"
)


func GetUserByID(repo repository.UserRepository, userID int) (UserInfo, error) {
    user, err := repo.FindByID(userID)
    if err != nil {
        return UserInfo{}, err
    }
    return UserInfo{
        ID:       user.ID,
        Email:    user.Email,
        Username: user.Username,
    }, nil
}