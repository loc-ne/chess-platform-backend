package db

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
    Client   *mongo.Client
    Database *mongo.Database
}

func ConnectMongoDB() (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    mongoURI := "mongodb+srv://nakrothnguyen127:zkpHsO8fNFuuS1kZ@cluster0.fqgz7hu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
    
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().
        ApplyURI(mongoURI).
        SetServerAPIOptions(serverAPI).
        SetMaxPoolSize(50).                    // Max 50 connections
        SetMinPoolSize(10).                    // Min 10 connections  
        SetMaxConnIdleTime(30 * time.Second).  // Idle timeout
        SetConnectTimeout(10 * time.Second).   // Connection timeout
        SetSocketTimeout(30 * time.Second).    // Socket timeout
        SetServerSelectionTimeout(5 * time.Second) // Server selection timeout

    client, err := mongo.Connect(ctx, opts)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer pingCancel()
    
    if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
        client.Disconnect(ctx)
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    log.Println("✅ Successfully connected to MongoDB!")
    
    database := client.Database("kangyoo")

    return &MongoDB{
        Client:   client,
        Database: database,
    }, nil
}

func (m *MongoDB) Disconnect() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := m.Client.Disconnect(ctx); err != nil {
        return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
    }
    
    log.Println("Disconnected from MongoDB")
    return nil
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
    return m.Database.Collection(name)
}