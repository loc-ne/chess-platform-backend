package handler

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/locne/ws-service/internal/usecase"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true 
    },
}

type JoinRoomMessage struct {
    Type     string `json:"type"`
    RoomID   string `json:"roomId"`
    UserID   int    `json:"userId"`
    Username string `json:"username"`
}

type JoinMatchmakingMessage struct {
    Type     string `json:"type"`
    UserID   int    `json:"userId"`
    Username string `json:"username"`
}

type GameActionMessage struct {
    Type     string `json:"type"`     // "gameAction" 
    RoomID   string `json:"roomId"`
    PlayerID int    `json:"playerId"`
    Action   string `json:"action"`   // "resign", "drawOffer", "drawAccept", "drawDecline"
    OfferID  string `json:"offerId,omitempty"` // For draw offers
}

func RegisterWebSocketRoutes(router *gin.Engine, rm *usecase.RoomManager) {
    router.GET("/ws", func(c *gin.Context) {
        handleWebSocket(c.Writer, c.Request, rm)
    })
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, rm *usecase.RoomManager) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()

    var client *usecase.Client
    var currentRoomID string

    conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    conn.SetPongHandler(func(string) error {
        conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    send := make(chan []byte, 256)

    go func() {
        ticker := time.NewTicker(54 * time.Second)
        defer ticker.Stop()
        defer conn.Close()

        for {
            select {
            case message, ok := <-send:
                conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
                if !ok {
                    conn.WriteMessage(websocket.CloseMessage, []byte{})
                    return
                }

                if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
                    log.Printf("Write error: %v", err)
                    return
                }

            case <-ticker.C:
                conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
                if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                    return
                }
            }
        }
    }()

    for {
        _, messageData, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }

        var rawMessage map[string]interface{}
        if err := json.Unmarshal(messageData, &rawMessage); err != nil {
            log.Printf("Invalid JSON: %v", err)
            continue
        }

        msgType, ok := rawMessage["type"].(string)
        if !ok {
            log.Printf("Missing message type")
            continue
        }

        switch msgType {
        case "joinMatchmaking":
            var joinMsg JoinMatchmakingMessage
            if err := json.Unmarshal(messageData, &joinMsg); err != nil {
                log.Printf("Invalid joinMatchmaking message: %v", err)
                continue
            }

            client = &usecase.Client{
                Conn:     conn,
                UserID:   joinMsg.UserID,
                Username: joinMsg.Username,
                Send:     send,
            }

            currentRoomID = "matchmaking"
            rm.JoinRoom(currentRoomID, client)
            
            log.Printf("User %d joined matchmaking", joinMsg.UserID)

        case "joinRoom":
            var joinMsg JoinRoomMessage
            if err := json.Unmarshal(messageData, &joinMsg); err != nil {
                log.Printf("Invalid joinRoom message: %v", err)
                continue
            }

            client = &usecase.Client{
                Conn:     conn,
                UserID:   joinMsg.UserID,
                Username: joinMsg.Username,
                Send:     send,
            }

            currentRoomID = joinMsg.RoomID
            rm.JoinRoom(currentRoomID, client)

            getStateMsg := usecase.MoveMessage{
                Type:     "getGameState",
                RoomID:   currentRoomID,
                PlayerID: joinMsg.UserID,
            }
        
            if err := rm.PublishMove(getStateMsg); err != nil {
                log.Printf("Failed to request game state: %v", err)
                rm.SendErrorToClient(currentRoomID, client.UserID, "Failed to get game state")
            }
        
            log.Printf("User %d joined room %s and requested game state", joinMsg.UserID, currentRoomID)

        case "move":
            if client == nil {
                log.Printf("Move without joining room")
                continue
            }

            var moveMsg usecase.MoveMessage
            if err := json.Unmarshal(messageData, &moveMsg); err != nil {
                log.Printf("Invalid move message: %v", err)
                continue
            }

            moveMsg.Type = "move"
            moveMsg.RoomID = currentRoomID
            moveMsg.PlayerID = client.UserID

            if err := rm.PublishMove(moveMsg); err != nil {
                log.Printf("Failed to publish move: %v", err)
                rm.SendErrorToClient(currentRoomID, client.UserID, "Failed to process move")
            }

        case "gameAction":
            if client == nil {
                log.Printf("Game action without joining room")
                continue
            }

            var actionMsg GameActionMessage
            if err := json.Unmarshal(messageData, &actionMsg); err != nil {
                log.Printf("Invalid game action message: %v", err)
                continue
            }

            // Convert to usecase.GameActionMessage
            gameActionMsg := usecase.GameActionMessage{
                Type:     "gameAction",
                RoomID:   currentRoomID,
                PlayerID: client.UserID,
                Action:   actionMsg.Action,
                OfferID:  actionMsg.OfferID,
            }

            if err := rm.PublishGameAction(gameActionMsg); err != nil {
                log.Printf("Failed to publish game action: %v", err)
                rm.SendErrorToClient(currentRoomID, client.UserID, "Failed to process game action")
            }

            log.Printf("User %d sent game action %s in room %s", client.UserID, actionMsg.Action, currentRoomID)

        case "leaveRoom":
            if client != nil && currentRoomID != "" {
                rm.LeaveRoom(currentRoomID, client.UserID)
            }
            return

        default:
            log.Printf("Unknown message type: %s", msgType)
        }
    }

    if client != nil && currentRoomID != "" {
        rm.LeaveRoom(currentRoomID, client.UserID)
    }
}