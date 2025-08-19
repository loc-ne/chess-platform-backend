package game

import (
    "fmt"
    "time"
    "github.com/locne/game-service/internal/usecase/engine"
    //"github.com/locne/game-service/internal/interface/repository"
    "github.com/locne/game-service/internal/entity"
    "strconv"
)

func (g *Game) MakeMove(playerID int, from, to engine.Position, gm *GameManager) error {
    g.mutex.Lock()
    defer g.mutex.Unlock()
    
    // 1. Check if spectator trying to move 
    if _, isSpectator := g.Spectators[playerID]; isSpectator {
        return fmt.Errorf("spectators cannot make moves")
    }
    
    // 2. Get and validate player
    player, exists := g.Players[playerID]
    if !exists {
        return fmt.Errorf("player not found in this game")
    }
    
    // 3. Check if player is online
    if !player.IsOnline {
        return fmt.Errorf("player is offline")
    }
    
    // 4. Validate turn
    if player.Color != g.GameState.ActiveColor {
        return fmt.Errorf("not your turn")
    }
    
    // 5. Update time before move
    g.updatePlayerTime()
    
    // 6. Check if player has time left
    if g.GameState.ActiveColor == "white" && g.WhiteTimeLeft <= 0 {
        g.handleTimeOut("white", gm)
        return fmt.Errorf("white time expired")
    }
    if g.GameState.ActiveColor == "black" && g.BlackTimeLeft <= 0 {
        g.handleTimeOut("black", gm)
        return fmt.Errorf("black time expired")
    }
    
    // 7. Validate move positions
    if err := g.validateMovePositions(from, to); err != nil {
        return err
    }

    // 8. Store state before move for notation
    gameBefore := g.GameState.Bitboards
    stateBefore := engine.GameState{
        ActiveColor:     g.GameState.ActiveColor,
        CastlingRights:  g.GameState.CastlingRights,
        EnPassantSquare: g.GameState.EnPassantSquare,
        MoveCount:       g.GameState.FullMoveNumber,
        HalfMoveClock:   g.GameState.HalfMoveClock,
    }
    
    // 9. Execute move using chess engine
    chessEngine := &engine.ChessEngine{}
    success := chessEngine.ExecuteServerMove(g.GameState, from, to)
    if !success {
        return fmt.Errorf("invalid move: from (%d,%d) to (%d,%d)", from.Row, from.Col, to.Row, to.Col)
    }

    // 10. Build notation
    notation := chessEngine.BuildNotation(
        gameBefore,
        stateBefore,
        g.GameState.Bitboards,
        engine.GameState{
            ActiveColor:     g.GameState.ActiveColor,
            CastlingRights:  g.GameState.CastlingRights,
            EnPassantSquare: g.GameState.EnPassantSquare,
            MoveCount:       g.GameState.FullMoveNumber,
            HalfMoveClock:   g.GameState.HalfMoveClock,
        },
        from, to,
    )
    g.addNotationToMoveHistory(notation)

    // 11. Post-move actions
    g.addTimeIncrement()
    g.LastMoveTime = time.Now()
    g.UpdatedAt = time.Now()
    
    g.switchTurn()

    // 12. Check if game ended
    finished, winner, reason := g.isGameFinished(gm)
    if finished {
        g.endGame(winner, reason, gm)
    }

    return nil
}

func (g *Game) switchTurn() {
    if g.GameState.ActiveColor == "white" {
        g.GameState.ActiveColor = "black"
    } else {
        g.GameState.ActiveColor = "white"
    }
}

func (g *Game) addNotationToMoveHistory(notation string) {
    isWhite := g.GameState.ActiveColor == "white"
    if isWhite {
        moveNum := len(g.GameState.MoveHistory) + 1
        g.GameState.MoveHistory = append(g.GameState.MoveHistory, engine.MoveNotation{
            MoveNumber: moveNum,
            White:      notation,
            Black:      "",
        })
    } else {
        if len(g.GameState.MoveHistory) == 0 {
            g.GameState.MoveHistory = append(g.GameState.MoveHistory, engine.MoveNotation{
                MoveNumber: 1,
                White:      "",
                Black:      notation,
            })
        } else {
            g.GameState.MoveHistory[len(g.GameState.MoveHistory)-1].Black = notation
        }
    }
}

func (g *Game) validateMovePositions(from, to engine.Position) error {
    if from.Row < 0 || from.Row > 7 || from.Col < 0 || from.Col > 7 {
        return fmt.Errorf("invalid 'from' position: row and col must be 0-7")
    }
    if to.Row < 0 || to.Row > 7 || to.Col < 0 || to.Col > 7 {
        return fmt.Errorf("invalid 'to' position: row and col must be 0-7")
    }
    if from.Row == to.Row && from.Col == to.Col {
        return fmt.Errorf("cannot move piece to the same position")
    }
    return nil
}

func (g *Game) updatePlayerTime() {
    if g.LastMoveTime.IsZero() || g.GameState == nil {
        return
    }
    
    elapsed := int(time.Since(g.LastMoveTime).Seconds())
    
    if g.GameState.ActiveColor == "white" {
        g.WhiteTimeLeft -= elapsed
        if g.WhiteTimeLeft < 0 {
            g.WhiteTimeLeft = 0
        }
    } else {
        g.BlackTimeLeft -= elapsed
        if g.BlackTimeLeft < 0 {
            g.BlackTimeLeft = 0
        }
    }
}

func (g *Game) addTimeIncrement() {
    if g.TimeControl == nil || g.TimeControl.Increment == 0 {
        return
    }
    
    if g.GameState.ActiveColor == "white" {
        g.WhiteTimeLeft += g.TimeControl.Increment
    } else {
        g.BlackTimeLeft += g.TimeControl.Increment
    }
}

func (g *Game) handleTimeOut(color string, gm *GameManager) {
    winner := "black"
    if color == "black" {
        winner = "white"
    }
    reason := "timeout"

    g.endGame(winner, reason, gm)
}

func (g *Game) isGameFinished(gm *GameManager) (bool, string, string) {
    if g.GameState == nil {
        return false, "", ""
    }
    
    chessEngine := &engine.ChessEngine{}
    finished, reason := chessEngine.IsGameOver(
        g.GameState.Bitboards,
        engine.GameState{
            ActiveColor:     g.GameState.ActiveColor,
            CastlingRights:  g.GameState.CastlingRights,
            EnPassantSquare: g.GameState.EnPassantSquare,
            MoveCount:       g.GameState.FullMoveNumber,
            HalfMoveClock:   g.GameState.HalfMoveClock,
        },
        g.GameState.PositionCounts,
        g.GameState.CurrentFen,
    )

    winner := "none"
    if finished {
        drawReasons := map[string]bool{
            "threefold repetition":   true,
            "stalemate":              true,
            "insufficient material":  true,
            "fifty move rule":        true,
        }

        if drawReasons[reason] {
            winner = "none"
        } else if g.GameState.ActiveColor == "black" {
            winner = "white"
        } else {
            winner = "black"
        }

    }
    
    return finished, winner, reason
}

func (g *Game) endGame(winner string, reason string, gm *GameManager) {
    whitePlayer := g.getPlayerByColor("white")
    blackPlayer := g.getPlayerByColor("black")
    
    result := "1/2-1/2" // draw
    winnerId := "none"
    
    if winner == "white" {
        result = "1-0"
        winnerId = strconv.Itoa(whitePlayer.ID)
    } else if winner == "black" {
        result = "0-1"
        winnerId = strconv.Itoa(blackPlayer.ID)
    }

    chessEngine := &engine.ChessEngine{}
    castlingStr := g.GameState.CastlingRights.ToFEN()
    enPassantStr := g.GameState.EnPassantSquare.ToFEN()
    
    game := entity.Game{
        GameID: g.ID,
        Players: struct {
            White entity.PlayerInfo `bson:"white"`
            Black entity.PlayerInfo `bson:"black"`
        }{
            White: entity.PlayerInfo{
                UserID:   whitePlayer.ID,
                Username: whitePlayer.Username,
                Elo:      whitePlayer.Rating,
            },
            Black: entity.PlayerInfo{
                UserID:   blackPlayer.ID,
                Username: blackPlayer.Username,
                Elo:      blackPlayer.Rating,
            },
        },
        Moves:         moveHistoryToNotationList(g.GameState.MoveHistory),
        Result:        result,
        CreatedAt:     g.CreatedAt,
        TimeControl:   strconv.Itoa(g.TimeControl.InitialTime/60) + "+" + strconv.Itoa(g.TimeControl.Increment),
        GameType:      g.TimeControl.Type,
        WinnerID:      winnerId,
        WhiteTimeLeft: g.WhiteTimeLeft,
        BlackTimeLeft: g.BlackTimeLeft,
        Reason:        reason,
        LastFen: chessEngine.BitboardToFEN(
            g.GameState.Bitboards,
            g.GameState.ActiveColor,
            castlingStr,
            enPassantStr,
            g.GameState.HalfMoveClock,
            g.GameState.FullMoveNumber,
        ),
    }
    fmt.Print(game)

    endUpdate := StateUpdateMessage{
        Type:   "gameEnd",
        RoomID: g.ID,
        Result: result,
        Winner: winner,
        Reason: reason,
    }
    gm.PublishStateUpdate(endUpdate)
    
    gm.RemoveGame(g.ID)

    gm.savePool.SaveGame(game)
}

func (g *Game) getPlayerByColor(color string) *Player {
    for _, player := range g.Players {
        if player.Color == color {
            return player
        }
    }
    return nil
}

func moveHistoryToNotationList(history []engine.MoveNotation) []string {
    var moves []string
    for _, m := range history {
        if m.White != "" {
            moves = append(moves, m.White)
        }
        if m.Black != "" {
            moves = append(moves, m.Black)
        }
    }
    return moves
}