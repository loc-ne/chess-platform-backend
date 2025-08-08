package usecase

import (
    "time"
    "crypto/rand"
    "encoding/hex"
    "github.com/golang-jwt/jwt/v5"
    "github.com/locne/auth-service/internal/entity"
	"github.com/locne/auth-service/internal/interface/repository"
	"os"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
    Sub      int    `json:"sub"`
    Email    string `json:"email"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateTokens(userRepo repository.UserRepository, user entity.User) (accessToken string, refreshToken string, err error) {
    claims := JWTClaims{
        Sub:      user.ID,
        Email:    user.Email,
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   string(rune(user.ID)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    accessToken, err = token.SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }

    b := make([]byte, 32)
    _, err = rand.Read(b)
    if err != nil {
        return "", "", err
    }
    refreshToken = hex.EncodeToString(b)

    err = userRepo.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

    return accessToken, refreshToken, nil
}