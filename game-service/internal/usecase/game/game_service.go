package game

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "sync"
    "github.com/go-redis/redis/v8"
    "github.com/locne/game-service/internal/usecase/engine"
    "github.com/locne/game-service/internal/interface/repository"
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

type DrawOffer struct {
    ID       string    `json:"id"`
    FromID   int       `json:"fromId"`
    ToID     int       `json:"toId"`
    CreatedAt time.Time `json:"createdAt"`
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
    DrawOffers    map[string]*DrawOffer  `json:"drawOffers"` // Active draw offers
    mutex         sync.RWMutex
}

type GameManager struct {
    redis    *redis.Client
    games    map[string]*Game
    mutex    sync.RWMutex
    ctx      context.Context
    savePool *GameSaveWorkerPool
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

type GameActionMessage struct {
    Type     string `json:"type"` // "gameAction"
    RoomID   string `json:"roomId"`
    PlayerID int    `json:"playerId"`
    Action   string `json:"action"` // "resign", "drawOffer", "drawAccept", "drawDecline"
    OfferID  string `json:"offerId,omitempty"` // For draw offers
}

type StateUpdateMessage struct {
    Type          string                  `json:"type"`
    RoomID        string                  `json:"roomId"`
    GameState     engine.ClientGameState  `json:"gameState,omitempty"`
    Player1       Player                  `json:"player1,omitempty"`
    Player2       Player                  `json:"player2,omitempty"`
    WhiteTimeLeft int                     `json:"whiteTimeLeft,omitempty"`
    BlackTimeLeft int                     `json:"blackTimeLeft,omitempty"`
    MoveHistory   []engine.MoveNotation   `json:"moveHistory,omitempty"`
    Error         string                  `json:"error,omitempty"`
    Result        string                  `json:"result,omitempty"`
    Reason        string                  `json:"reason,omitempty"`
    Winner        string                  `json:"winner,omitempty"`
    OfferID       string                  `json:"offerId,omitempty"` // For draw offers
    OfferFrom     int                     `json:"offerFrom,omitempty"` // Player ID who made the offer
    TargetPlayerID *int                   `json:"targetPlayerId,omitempty"` // For targeted messages
}

func NewGameManager(redis *redis.Client, ctx context.Context, repo repository.GameRepository) *GameManager {
    return &GameManager{
        redis:    redis,
        games:    make(map[string]*Game),
        savePool: NewGameSaveWorkerPool(repo, 3),
        ctx:      ctx,
    }
}

func (gm *GameManager) ListenChannels() {
    // Start move listener
    go gm.ListenMoves()
    
    // Start game action listener
    go gm.ListenGameActions()
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

func (gm *GameManager) ListenGameActions() {
    pubsub := gm.redis.Subscribe(gm.ctx, "game_action")
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        var actionMsg GameActionMessage
        if err := json.Unmarshal([]byte(msg.Payload), &actionMsg); err != nil {
            continue
        }
        gm.ProcessGameAction(actionMsg)
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
        return
    }

    from := engine.Position{Row: moveMsg.FromRow, Col: moveMsg.FromCol}
    to := engine.Position{Row: moveMsg.ToRow, Col: moveMsg.ToCol}
    
    err := game.MakeMove(moveMsg.PlayerID, from, to, gm)
    if err != nil {
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
        MoveHistory: game.GameState.MoveHistory,
    }
    
    gm.PublishStateUpdate(stateUpdate)
}

func (gm *GameManager) ProcessGameAction(actionMsg GameActionMessage) {
    gm.mutex.RLock()
    game, exists := gm.games[actionMsg.RoomID]
    gm.mutex.RUnlock()
    
    if !exists {
        return
    }

    switch actionMsg.Action {
    case "resign":
        gm.handleResign(game, actionMsg.PlayerID)
    case "drawOffer":
        gm.handleDrawOffer(game, actionMsg.PlayerID)
    case "drawAccept":
        gm.handleDrawAccept(game, actionMsg.PlayerID, actionMsg.OfferID)
    case "drawDecline":
        gm.handleDrawDecline(game, actionMsg.PlayerID, actionMsg.OfferID)
    default:
        gm.PublishError(actionMsg.RoomID, "Unknown action: "+actionMsg.Action)
    }
}

func (gm *GameManager) handleResign(game *Game, playerID int) {
    game.mutex.Lock()
    var winner string
    for _, player := range game.Players {
        if player.ID != playerID {
            winner = player.Color
            break
        }
    }

    game.mutex.Unlock() 
    
    
    game.endGame(winner, "resignation", gm)
    
}

func (gm *GameManager) handleDrawOffer(game *Game, playerID int) {
    game.mutex.Lock()
    defer game.mutex.Unlock()
    
    // Generate offer ID
    offerID := fmt.Sprintf("%s_%d_%d", game.ID, playerID, time.Now().Unix())
    
    // Find opponent
    var opponentID int
    for _, player := range game.Players {
        if player.ID != playerID {
            opponentID = player.ID
            break
        }
    }
    
    // Create draw offer
    offer := &DrawOffer{
        ID:        offerID,
        FromID:    playerID,
        ToID:      opponentID,
        CreatedAt: time.Now(),
    }
    
    // Initialize DrawOffers map if nil
    if game.DrawOffers == nil {
        game.DrawOffers = make(map[string]*DrawOffer)
    }
    
    // Store offer
    game.DrawOffers[offerID] = offer
    
    // Notify opponent only (targeted message)
    update := StateUpdateMessage{
        Type:           "drawOffer",
        RoomID:         game.ID,
        OfferID:        offerID,
        OfferFrom:      playerID,
        TargetPlayerID: &opponentID, 
    }
    
    gm.PublishStateUpdate(update)
}

func (gm *GameManager) handleDrawAccept(game *Game, playerID int, offerID string) {
    game.mutex.Lock()
    
    // Check if offer exists
    offer, exists := game.DrawOffers[offerID]
    if !exists {
        game.mutex.Unlock() 
        return
    }
    
    // Check if player is the recipient of the offer
    if offer.ToID != playerID {
        game.mutex.Unlock() 
        return
    }
    
    game.DrawOffers = make(map[string]*DrawOffer)
    
    game.mutex.Unlock() 

    game.endGame("", "draw by agreement", gm) 
}

func (gm *GameManager) handleDrawDecline(game *Game, playerID int, offerID string) {
    game.mutex.Lock()
    defer game.mutex.Unlock()
    
    // Check if offer exists
    offer, exists := game.DrawOffers[offerID]
    if !exists {
        return
    }
    
    // Check if player is the recipient of the offer
    if offer.ToID != playerID {
        return
    }
    
    // Remove offer
    delete(game.DrawOffers, offerID)
    
    // Notify that offer was declined
    update := StateUpdateMessage{
        Type:    "drawDeclined",
        RoomID:  game.ID,
        OfferID: offerID,
    }
    
    gm.PublishStateUpdate(update)
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