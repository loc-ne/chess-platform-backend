package usecase

import (
    "github.com/locne/auth-service/internal/entity"
	"github.com/locne/auth-service/internal/interface/repository"
	"golang.org/x/crypto/bcrypt"
	"regexp"
    "fmt"
)

type UserInfo struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
}

func isEmail(text string) bool {
    emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+$`
    re := regexp.MustCompile(emailRegex)
    return re.MatchString(text)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func validateUser(userRepo repository.UserRepository, username, password string) (entity.User, error) {
    var user entity.User
    var err error

    if isEmail(username) {
        user, err = userRepo.FindByEmail(username)
    } else {
        user, err = userRepo.FindByUsername(username)
    }
    if err != nil {
        return user, err
    }
    if user.ID == 0 {
        return user, fmt.Errorf("invalid credentials")
    }
    if !user.IsActive {
        return user, fmt.Errorf("account is disabled")
    }
    if !CheckPasswordHash(password, user.Password) {
        return user, fmt.Errorf("invalid credentials")
    }
    return user, nil
}

func Login(userRepo repository.UserRepository, username, password string) (UserInfo, string, string, error) {
    user, err := validateUser(userRepo, username, password)
    if err != nil {
        return UserInfo{}, "", "", err
    }
    accessToken, refreshToken, err := GenerateTokens(userRepo, user)
    if err != nil {
        return UserInfo{}, "", "", err
    }
    userInfo := UserInfo{
        ID:       user.ID,
        Email:    user.Email,
        Username: user.Username,
    }
    return userInfo, accessToken, refreshToken, nil
}

func Register(userRepo repository.UserRepository, user entity.User) (UserInfo, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return user, err
    }
    user.Password = string(hashedPassword)
    err = userRepo.Create(user)

    if err != nil {
        return user, err
    }

    userInfo := UserInfo{
        ID:       user.ID,
        Email:    user.Email,
        Username: user.Username,
    }

    return userInfo, nil
}