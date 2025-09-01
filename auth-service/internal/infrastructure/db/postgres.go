package db

import (
    "fmt"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectPostgres() (*gorm.DB, error) {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
    }

    maxOpenConns := 10
    maxIdleConns := 5
    connMaxLifetime := 30 * time.Minute
    retryAttempts := 5
    retryDelay := 2 * time.Second

    var db *gorm.DB
    var err error

    for attempt := 1; attempt <= retryAttempts; attempt++ {
        db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
        if err == nil {
            sqlDB, err := db.DB()
            if err != nil {
                return nil, err
            }
            sqlDB.SetMaxOpenConns(maxOpenConns)
            sqlDB.SetMaxIdleConns(maxIdleConns)
            sqlDB.SetConnMaxLifetime(connMaxLifetime)

            if err := sqlDB.Ping(); err != nil {
                return nil, fmt.Errorf("database ping failed: %w", err)
            }

            return db, nil
        }

        fmt.Printf("DB connect attempt %d failed: %v\n", attempt, err)
        time.Sleep(retryDelay)
    }

    return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", retryAttempts, err)
}

func CloseDB(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
