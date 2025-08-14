package usecase

import (
    "context"
    "encoding/json"
    "log"
    "sync"
    "github.com/go-redis/redis/v8"
    "github.com/gorilla/websocket"
)

type Client struct {
    Conn     *websocket.Conn
    UserID   int
    Username string
    Send     chan []byte
}

type Room struct {
    ID      string
    Clients map[int]*Client 
    mutex   sync.RWMutex
}

type RoomManager struct {
    redis *redis.Client
    rooms map[string]*Room
    mutex sync.RWMutex
    ctx   context.Context
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

type Player struct {
    ID       int    `json:"userId"`       
    Username string `json:"username"`
    Color    string `json:"color"`       
    Rating   int    `json:"rating,omitempty"`   
    IsOnline bool   `json:"isOnline"`   
}

type Position struct {
    Row int `json:"row"`
    Col int `json:"col"`
}

type CastlingRights struct {
    WhiteKingSide  bool `json:"whiteKingSide"`
    WhiteQueenSide bool `json:"whiteQueenSide"`
    BlackKingSide  bool `json:"blackKingSide"`
    BlackQueenSide bool `json:"blackQueenSide"`
}

type BitboardGame struct {
    WhitePawns   uint64 `json:"WhitePawns,string"`
    WhiteRooks   uint64 `json:"WhiteRooks,string"`
    WhiteKnights uint64 `json:"WhiteKnights,string"`
    WhiteBishops uint64 `json:"WhiteBishops,string"`
    WhiteQueens  uint64 `json:"WhiteQueens,string"`
    WhiteKing    uint64 `json:"WhiteKing,string"`

    BlackPawns   uint64 `json:"BlackPawns,string"`
    BlackRooks   uint64 `json:"BlackRooks,string"`
    BlackKnights uint64 `json:"BlackKnights,string"`
    BlackBishops uint64 `json:"BlackBishops,string"`
    BlackQueens  uint64 `json:"BlackQueens,string"`
    BlackKing    uint64 `json:"BlackKing,string"`
}

type ClientGameState struct {
    CurrentFen      string          `json:"currentFen"`
    Bitboards       BitboardGame    `json:"bitboards"`
    ActiveColor     string          `json:"activeColor"`
    CastlingRights  CastlingRights  `json:"castlingRights"`
    EnPassantSquare *Position       `json:"enPassantSquare"`
}

type StateUpdateMessage struct {
    Type          string          `json:"type"`
    RoomID        string          `json:"roomId"`
    GameState     ClientGameState `json:"gameState,omitempty"`
    Player1       Player          `json:"player1,omitempty"`
    Player2       Player          `json:"player2,omitempty"`
    WhiteTimeLeft int             `json:"whiteTimeLeft,omitempty"`
    BlackTimeLeft int             `json:"blackTimeLeft,omitempty"`
    Error         string          `json:"error,omitempty"`
    Result        string          `json:"result,omitempty"`
    Reason        string          `json:"reason,omitempty"`
}

func NewRoomManager(redisClient *redis.Client) *RoomManager {
    return &RoomManager{
        redis: redisClient,
        rooms: make(map[string]*Room),
        ctx:   context.Background(),
    }
}

func (rm *RoomManager) ListenStateUpdates() {
    pubsub := rm.redis.Subscribe(rm.ctx, "move_out")
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        var stateUpdate StateUpdateMessage
        if err := json.Unmarshal([]byte(msg.Payload), &stateUpdate); err != nil {
            log.Printf("Error unmarshaling state update: %v", err)
            continue
        }
        
        if stateUpdate.Type == "matchFound" {
            rm.handleMatchFound(stateUpdate)
        } else {
            rm.BroadcastToRoom(stateUpdate.RoomID, stateUpdate)
        }
    }
}

func (rm *RoomManager) handleMatchFound(matchFound StateUpdateMessage) {
    rm.mutex.RLock()
    matchmakingRoom, exists := rm.rooms["matchmaking"]
    rm.mutex.RUnlock()

    if !exists {
        log.Printf("No matchmaking room found")
        return
    }

    targetUserIDs := []int{
        matchFound.Player1.ID,
        matchFound.Player2.ID,
    }

    data, err := json.Marshal(matchFound)
    if err != nil {
        log.Printf("Error marshaling matchFound message: %v", err)
        return
    }

    matchmakingRoom.mutex.RLock()
    for _, userID := range targetUserIDs {
        if client, exists := matchmakingRoom.Clients[userID]; exists {
            select {
            case client.Send <- data:
                log.Printf("Sent matchFound to user %d", userID)
            default:
                log.Printf("Failed to send matchFound to user %d", userID)
            }
        }
    }
    matchmakingRoom.mutex.RUnlock()
}

func (rm *RoomManager) PublishMove(moveMsg MoveMessage) error {
    data, err := json.Marshal(moveMsg)
    if err != nil {
        return err
    }

    return rm.redis.Publish(rm.ctx, "move_in", data).Err()
}

func (rm *RoomManager) BroadcastToRoom(roomID string, message interface{}) {
    rm.mutex.RLock()
    room, exists := rm.rooms[roomID]
    rm.mutex.RUnlock()

    if !exists {
        log.Printf("Room %s not found for broadcast", roomID)
        return
    }

    data, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshaling broadcast message: %v", err)
        return
    }

    room.mutex.RLock()
    for _, client := range room.Clients {
        select {
        case client.Send <- data:
        default:
            log.Printf("Client %d send channel full, skipping", client.UserID)
        }
    }
    room.mutex.RUnlock()
}


func (rm *RoomManager) JoinRoom(roomID string, client *Client) {
    rm.mutex.Lock()
    room, exists := rm.rooms[roomID]
    if !exists {
        room = &Room{
            ID:      roomID,
            Clients: make(map[int]*Client),
        }
        rm.rooms[roomID] = room
    }
    rm.mutex.Unlock()

    room.mutex.Lock()
    room.Clients[client.UserID] = client
    room.mutex.Unlock()

    log.Printf("User %d (%s) joined room %s", client.UserID, client.Username, roomID)
}

func (rm *RoomManager) LeaveRoom(roomID string, userID int) {
    rm.mutex.RLock()
    room, exists := rm.rooms[roomID]
    rm.mutex.RUnlock()

    if !exists {
        return
    }

    room.mutex.Lock()
    if client, exists := room.Clients[userID]; exists {
        close(client.Send)
        delete(room.Clients, userID)
        log.Printf("User %d left room %s", userID, roomID)
    }
    room.mutex.Unlock()

    room.mutex.RLock()
    isEmpty := len(room.Clients) == 0
    room.mutex.RUnlock()

    if isEmpty {
        rm.mutex.Lock()
        delete(rm.rooms, roomID)
        rm.mutex.Unlock()
        log.Printf("Room %s deleted (empty)", roomID)
    }
}


func (rm *RoomManager) SendErrorToClient(roomID string, userID int, errorMsg string) {
    rm.mutex.RLock()
    room, exists := rm.rooms[roomID]
    rm.mutex.RUnlock()

    if !exists {
        return
    }

    room.mutex.RLock()
    client, exists := room.Clients[userID]
    room.mutex.RUnlock()

    if !exists {
        return
    }

    errorResponse := map[string]interface{}{
        "type":  "error",
        "error": errorMsg,
    }

    data, _ := json.Marshal(errorResponse)
    select {
    case client.Send <- data:
    default:
        log.Printf("Failed to send error to client %d", userID)
    }
}