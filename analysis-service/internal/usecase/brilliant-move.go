package usecase

import (
    "github.com/notnil/chess"
    "fmt"
)

type PieceValue struct {
    Pawn   int
    Knight int
    Bishop int
    Rook   int
    Queen  int
}

var StandardPieceValues = PieceValue{
    Pawn:   1,
    Knight: 3,
    Bishop: 3,
    Rook:   5,
    Queen:  9,
}

type MoveInfo struct {
    Move          *chess.Move
    CapturedPiece chess.PieceType
    MovingPiece   chess.PieceType
    IsCapture     bool
    IsSacrifice   bool
    MaterialLoss  int
}

func getMoveInfo(game *chess.Game, moveStr string) (*MoveInfo, error) {
    move, err := findValidMove(game, moveStr)
    if err != nil {
        return nil, err
    }
    
    position := game.Position()
    movingPiece := position.Board().Piece(move.S1()).Type()
    
    info := &MoveInfo{
        Move:        move,
        MovingPiece: movingPiece,
        IsCapture:   move.HasTag(chess.Capture),
    }
    
    if info.IsCapture {
        info.CapturedPiece = position.Board().Piece(move.S2()).Type()
    }
    
    return info, nil
}

func isBrilliantMove(inputMoves []string, moveIndex int) bool {
    if moveIndex <= 0 {
        return false
        
    }
    
    // Create game state before the move
    game := chess.NewGame()
    for i := 0; i < moveIndex; i++ {
        move, err := findValidMove(game, inputMoves[i])
        if err != nil {
            return false
        }
        game.Move(move)
    }
    
    currentMove := inputMoves[moveIndex]
    moveInfo, err := getMoveInfo(game, currentMove)
    if err != nil {
        fmt.Println(err)
        return false
    }

    // Check if it's a sacrifice of valuable piece
    if !isSacrificeMove(moveInfo) {
        return false
    }
    
    // Execute the move to check opponent's response
    game.Move(moveInfo.Move)
    // Check if opponent can capture with lower value piece
    return canOpponentCaptureWithLowerValue(game, moveInfo)
}

func isSacrificeMove(info *MoveInfo) bool {
    // Must be moving a valuable piece (Knight, Bishop, Rook, Queen)
    movingValue := getPieceValue(info.MovingPiece)
    if movingValue <= StandardPieceValues.Pawn {
        return false
    }
    
    // If it's a capture, check if we're trading up or equal
    if info.IsCapture {
        capturedValue := getPieceValue(info.CapturedPiece)
        // Only sacrifice if we're giving up more than we take
        if movingValue == capturedValue {
            return false // Equal trade - not a sacrifice
        }
        return movingValue > capturedValue
    }
    
    // Non-capture move with valuable piece - check if piece is under attack
    return true
}

func canOpponentCaptureWithLowerValue(game *chess.Game, sacrificeInfo *MoveInfo) bool {
    validMoves := game.ValidMoves()
    myColor := game.Position().Turn().Other() 
    
    for _, move := range validMoves {
        if move.HasTag(chess.Capture) {
            capturingPiece := game.Position().Board().Piece(move.S1()).Type()
            capturingValue := getPieceValue(capturingPiece)
            
            capturedPiece := game.Position().Board().Piece(move.S2()).Type()
            capturedValue := getPieceValue(capturedPiece)
            
            capturedPieceColor := game.Position().Board().Piece(move.S2()).Color()
            if capturedPieceColor != myColor {
                continue 
            }
            
            if capturedValue >= StandardPieceValues.Knight && capturingValue < capturedValue {
                return true
            }
        }
    }
    
    return false
}

func getPieceValue(pieceType chess.PieceType) int {
    switch pieceType {
    case chess.Pawn:
        return StandardPieceValues.Pawn
    case chess.Knight:
        return StandardPieceValues.Knight
    case chess.Bishop:
        return StandardPieceValues.Bishop
    case chess.Rook:
        return StandardPieceValues.Rook
    case chess.Queen:
        return StandardPieceValues.Queen
    case chess.King:
        return 0 
    default:
        return 0
    }
}

func isPowerfulPiece(pieceType chess.PieceType) bool {
    return pieceType == chess.Knight || 
           pieceType == chess.Bishop || 
           pieceType == chess.Rook || 
           pieceType == chess.Queen
}
