package db

import (
    "fmt"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
    connStr := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    fmt.Println("Connected to Postgres with GORM")
    return db, nil
}