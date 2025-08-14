package game

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "sync"
    "github.com/go-redis/redis/v8"
    "github.com/locne/game-service/internal/usecase/engine"
)

type Player struct {
    ID       int    `json:"userId"`       
    Username string `json:"username"`
    Color    string `json:"color"`       
    Rating   int    `json:"rating,omitempty"`   
    IsOnline bool   `json:"isOnline"`   
}

type TimeControl struct {
    Type        string `json:"type"`        
    InitialTime int    `json:"initialTime"` 
    Increment   int    `json:"increment"`   
}

type Game struct {
    ID            string                 `json:"id"`
    Players       map[int]*Player        `json:"players"`
    Spectators    map[int]*Player        `json:"spectators"`
    GameState     *engine.ServerGameState `json:"gameState"`
    TimeControl   *TimeControl           `json:"timeControl"`
    WhiteTimeLeft int                    `json:"whiteTimeLeft"`
    BlackTimeLeft int                    `json:"blackTimeLeft"`
    LastMoveTime  time.Time              `json:"lastMoveTime"`
    CreatedAt     time.Time              `json:"createdAt"`
    UpdatedAt     time.Time              `json:"updatedAt"`
    mutex         sync.RWMutex
}

type GameManager struct {
    redis   *redis.Client
    games   map[string]*Game
    mutex   sync.RWMutex
    ctx     context.Context
}

type MoveMessage struct {
    Type      string `json:"type"`
    RoomID    string `json:"roomId"`
    PlayerID  int    `json:"playerId"`
    FromRow   int    `json:"fromRow"`
    FromCol   int    `json:"fromCol"`
    ToRow     int    `json:"toRow"`
    ToCol     int    `json:"toCol"`
    Promotion string `json:"promotion,omitempty"`
}

type StateUpdateMessage struct {
    Type          string                  `json:"type"`
    RoomID        string                  `json:"roomId"`
    GameState     engine.ClientGameState  `json:"gameState,omitempty"`
    Player1       Player                  `json:"player1,omitempty"`
    Player2       Player                  `json:"player2,omitempty"`
    WhiteTimeLeft int                     `json:"whiteTimeLeft,omitempty"`
    BlackTimeLeft int                     `json:"blackTimeLeft,omitempty"`
    Error         string                  `json:"error,omitempty"`
    Result        string                  `json:"result,omitempty"`
    Reason        string                  `json:"reason,omitempty"`
}

func NewGameManager(redis *redis.Client, ctx context.Context) *GameManager {
    return &GameManager{
        redis: redis,
        games: make(map[string]*Game),
        ctx:   ctx,
    }
}

func (gm *GameManager) ListenMoves() {
    pubsub := gm.redis.Subscribe(gm.ctx, "move_in")
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        var moveMsg MoveMessage
        if err := json.Unmarshal([]byte(msg.Payload), &moveMsg); err != nil {
            continue
        }
        gm.ProcessMove(moveMsg)
    }
}

func (gm *GameManager) ProcessMove(moveMsg MoveMessage) {
    if moveMsg.Type == "getGameState" {
        gameState, err := gm.GetGameState(moveMsg.RoomID)
        if err != nil {
            gm.PublishError(moveMsg.RoomID, err.Error())
            return
        }
        gm.PublishStateUpdate(*gameState)
        return
    }
    
    gm.mutex.RLock()
    game, exists := gm.games[moveMsg.RoomID]
    gm.mutex.RUnlock()
    
    if !exists {
        gm.PublishError(moveMsg.RoomID, "Game not found")
        return
    }

    from := engine.Position{Row: moveMsg.FromRow, Col: moveMsg.FromCol}
    to := engine.Position{Row: moveMsg.ToRow, Col: moveMsg.ToCol}
    
    err := game.MakeMove(moveMsg.PlayerID, from, to, gm)
    if err != nil {
        gm.PublishError(moveMsg.RoomID, err.Error())
        return
    }

    stateUpdate := StateUpdateMessage{
        Type:          "gameUpdate",
        RoomID:        moveMsg.RoomID,
        GameState:     engine.ClientGameState{
            CurrentFen:     game.GameState.CurrentFen,
            Bitboards:      game.GameState.Bitboards,
            ActiveColor:    game.GameState.ActiveColor,
            CastlingRights: game.GameState.CastlingRights,
            EnPassantSquare: game.GameState.EnPassantSquare,
        },
        WhiteTimeLeft: game.WhiteTimeLeft,
        BlackTimeLeft: game.BlackTimeLeft,
    }
    
    gm.PublishStateUpdate(stateUpdate)
}

func (gm *GameManager) GetGameState(gameID string) (*StateUpdateMessage, error) {
    gm.mutex.RLock()
    game, exists := gm.games[gameID]
    gm.mutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("game not found")
    }
    
    if game.GameState == nil {
        return nil, fmt.Errorf("game state is invalid")
    }
    
    var player1, player2 *Player
    for _, player := range game.Players {
        if player.Color == "white" {
            player1 = player
        } else if player.Color == "black" {
            player2 = player
        }
    }
    
    gameStateMsg := &StateUpdateMessage{
        Type:   "gameState",
        RoomID: gameID,
        GameState: engine.ClientGameState{
            CurrentFen:      game.GameState.CurrentFen,
            Bitboards:       game.GameState.Bitboards,
            ActiveColor:     game.GameState.ActiveColor,
            CastlingRights:  game.GameState.CastlingRights,
            EnPassantSquare: game.GameState.EnPassantSquare,
        },
        Player1:       *player1,
        Player2:       *player2,
        WhiteTimeLeft: game.WhiteTimeLeft,
        BlackTimeLeft: game.BlackTimeLeft,
    }
    
    return gameStateMsg, nil
}

func (gm *GameManager) PublishStateUpdate(update StateUpdateMessage) {
    data, _ := json.Marshal(update)
    gm.redis.Publish(gm.ctx, "move_out", data)
}

func (gm *GameManager) PublishError(roomID, errorMsg string) {
    errorUpdate := StateUpdateMessage{
        Type:   "error", 
        RoomID: roomID,
        Error:  errorMsg,
    }
    data, _ := json.Marshal(errorUpdate)
    gm.redis.Publish(gm.ctx, "move_out", data)
}

func (gm *GameManager) AddGame(game *Game) {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()
    gm.games[game.ID] = game
}

func (gm *GameManager) RemoveGame(gameID string) {
    gm.mutex.Lock()
    defer gm.mutex.Unlock()
    delete(gm.games, gameID)
}