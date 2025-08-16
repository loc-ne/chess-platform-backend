package repository

import (
    "context"
    "errors"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "github.com/locne/game-service/internal/entity"
    "fmt"
)	

type GameRepository interface {
    GetByID(ctx context.Context, id string) (entity.Game, error)
    GetAll(ctx context.Context) ([]entity.Game, error)
    Create(ctx context.Context, game entity.Game) error
    Update(ctx context.Context, game entity.Game) error
    Delete(ctx context.Context, id string) error
}

type mongoGameRepo struct {
    collection *mongo.Collection
}

func NewMongoGameRepo(db *mongo.Database) GameRepository {
    return &mongoGameRepo{
        collection: db.Collection("games"),
    }
}

func (r *mongoGameRepo) GetByID(ctx context.Context, gameId string) (entity.Game, error) {
    var game entity.Game
    filter := bson.M{"gameId": gameId}
    err := r.collection.FindOne(ctx, filter).Decode(&game)
    if err != nil {
        fmt.Printf("GetByID error: %v (gameId=%s)\n", err, gameId)
        return game, err
    }
    return game, nil
}

func (r *mongoGameRepo) GetAll(ctx context.Context) ([]entity.Game, error) {
    var games []entity.Game
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var game entity.Game
        if err := cursor.Decode(&game); err != nil {
            return nil, err
        }
        games = append(games, game)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return games, nil
}

func (r *mongoGameRepo) Create(ctx context.Context, game entity.Game) error {
    _, err := r.collection.InsertOne(ctx, game)
    return err
}

func (r *mongoGameRepo) Update(ctx context.Context, game entity.Game) error {
    if game.GameID == "" {
        return errors.New("missing gameId")
    }
    filter := bson.M{"gameId": game.GameID}
    update := bson.M{"$set": game}
    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

func (r *mongoGameRepo) Delete(ctx context.Context, gameId string) error {
    filter := bson.M{"gameId": gameId}
    _, err := r.collection.DeleteOne(ctx, filter)
    return err
}