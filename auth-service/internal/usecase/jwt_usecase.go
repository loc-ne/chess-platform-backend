package usecase

import (
    "time"
    "crypto/rand"
    "encoding/hex"
    "github.com/golang-jwt/jwt/v5"
    "github.com/locne/auth-service/internal/entity"
	"github.com/locne/auth-service/internal/interface/repository"
	"os"
    "fmt"
    "strings"
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

func RefreshToken(userRepo repository.UserRepository, refreshToken string) (accessToken string, newRefreshToken string, err error) {
    user, err := userRepo.FindByRefreshToken(refreshToken)
    if err != nil {
        return "", "", fmt.Errorf("invalid refresh token")
    }
    if !user.IsActive {
        return "", "", fmt.Errorf("user is not active")
    }
    return GenerateTokens(userRepo, user)
}

func ValidateToken(tokenString string) (int, error) {
    parts := strings.Split(tokenString, ".")
    if len(parts) != 3 {
        return 0, fmt.Errorf("invalid JWT format")
    }

    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })
    if err != nil {
        return 0, err
    }
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        if claims.ExpiresAt.After(time.Now()) {
            return claims.Sub, nil
        }
        return 0, fmt.Errorf("token expired")
    }
    return 0, fmt.Errorf("invalid token")
}