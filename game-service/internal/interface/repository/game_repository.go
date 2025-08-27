package repository

import (
    "context"
    "fmt"
    "log"
    
    "github.com/locne/game-service/internal/entity"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type GameRepository interface {
    SaveGame(ctx context.Context, game entity.Game) error
    SaveGamesBatch(ctx context.Context, games []entity.Game) error
    GetGameByID(ctx context.Context, gameID string) (*entity.Game, error)
}

type mongoGameRepository struct {
    collection *mongo.Collection
}

func NewGameRepository(db *mongo.Database) GameRepository {
    return &mongoGameRepository{
        collection: db.Collection("games"),
    }
}

func (r *mongoGameRepository) SaveGamesBatch(ctx context.Context, games []entity.Game) error {
    if len(games) == 0 {
        return nil
    }
    
    // Convert to interface slice for MongoDB
    documents := make([]interface{}, len(games))
    for i, game := range games {
        documents[i] = game
    }
    
    // Batch insert với ordered=false cho performance
    opts := options.InsertMany().SetOrdered(false)
    
    result, err := r.collection.InsertMany(ctx, documents, opts)
    if err != nil {
        return fmt.Errorf("batch insert failed: %w", err)
    }
    
    log.Printf("Batch inserted %d games, IDs: %v", 
        len(result.InsertedIDs), result.InsertedIDs)
    
    return nil
}

func (r *mongoGameRepository) SaveGame(ctx context.Context, game entity.Game) error {
    _, err := r.collection.InsertOne(ctx, game)
    if err != nil {
        return fmt.Errorf("failed to save game: %w", err)
    }
    return nil
}

func (r *mongoGameRepository) GetGameByID(ctx context.Context, gameID string) (*entity.Game, error) {
    var game entity.Game
    
    filter := bson.M{"gameId": gameID}
    err := r.collection.FindOne(ctx, filter).Decode(&game)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("game not found: %s", gameID)
        }
        return nil, fmt.Errorf("failed to get game: %w", err)
    }
    
    return &game, nil
}