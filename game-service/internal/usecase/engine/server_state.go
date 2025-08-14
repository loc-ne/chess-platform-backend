package engine

import ("fmt")

// Create initial server game state
func (ce *ChessEngine) CreateServerGameState() *ServerGameState {
    initialFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    
    return &ServerGameState{
        CurrentFen:  initialFen,
        Bitboards:   ce.CreateBitboardGame(),
        ActiveColor: "white",
        CastlingRights: CastlingRights{
            WhiteKingSide:  true,
            WhiteQueenSide: true,
            BlackKingSide:  true,
            BlackQueenSide: true,
        },
        EnPassantSquare: nil,
        //MoveHistory:     []MoveNotation,
        FullMoveNumber:  1,
        HalfMoveClock:   0,
        PositionCounts: map[string]int{
            initialFen: 1,
        },
        MaterialCount: map[string]MaterialCount{
            "white": {
                Pawns:   8,
                Knights: 2,
                Bishops: 2,
                Rooks:   2,
                Queens:  1,
            },
            "black": {
                Pawns:   8,
                Knights: 2,
                Bishops: 2,
                Rooks:   2,
                Queens:  1,
            },
        },
    }
}

// Update material count after move
func (ce *ChessEngine) UpdateMaterialCount(state *ServerGameState, capturedPiece *Piece) {
    if capturedPiece == nil {
        return
    }
    
    materialCount := state.MaterialCount[capturedPiece.Color]
    
    switch capturedPiece.Type {
    case "pawn":
        materialCount.Pawns--
    case "knight":
        materialCount.Knights--
    case "bishop":
        materialCount.Bishops--
    case "rook":
        materialCount.Rooks--
    case "queen":
        materialCount.Queens--
    }
    
    state.MaterialCount[capturedPiece.Color] = materialCount
}

// Update position counts for threefold repetition
func (ce *ChessEngine) UpdatePositionCounts(state *ServerGameState, newFen string) {
    state.PositionCounts[newFen]++
    state.CurrentFen = newFen
}

// Add move to history
// func (ce *ChessEngine) AddMoveToHistory(state *ServerGameState, move MoveNotation) {
//     state.MoveHistory = append(state.MoveHistory, move)
// }

// Convert move to algebraic notation
// func (ce *ChessEngine) MoveToAlgebraic(game BitboardGame, from Position, to Position) MoveNotation {
//     piece := ce.GetPieceAt(game, to)
//     if piece == nil {
//         return ""
//     }
    
//     fromSquare := ce.PositionToAlgebraic(from)
//     toSquare := ce.PositionToAlgebraic(to)
    
//     // Basic format: e2e4, Ng1f3, etc.
//     pieceSymbol := ""
//     switch piece.Type {
//     case "king":
//         pieceSymbol = "K"
//     case "queen":
//         pieceSymbol = "Q"
//     case "rook":
//         pieceSymbol = "R"
//     case "bishop":
//         pieceSymbol = "B"
//     case "knight":
//         pieceSymbol = "N"
//     case "pawn":
//         pieceSymbol = "" // Pawns don't have a symbol
//     }
    
//     return pieceSymbol + fromSquare + toSquare
// }

// Complete server move execution
func (ce *ChessEngine) BuildNotation(gameBefore BitboardGame, stateBefore GameState, gameAfter BitboardGame, stateAfter GameState, from, to Position) string {
    piece := ce.GetPieceAt(gameBefore, from)
    if piece == nil {
        return ""
    }

    // Ký hiệu quân cờ
    pieceSymbol := ""
    switch piece.Type {
    case "king":
        pieceSymbol = "K"
    case "queen":
        pieceSymbol = "Q"
    case "rook":
        pieceSymbol = "R"
    case "bishop":
        pieceSymbol = "B"
    case "knight":
        pieceSymbol = "N"
    case "pawn":
        pieceSymbol = ""
    }

    // Nhập thành
    if piece.Type == "king" && abs(to.Col-from.Col) == 2 {
        if to.Col == 6 {
            return "O-O"
        }
        if to.Col == 2 {
            return "O-O-O"
        }
    }

    // Phong cấp
    promotion := ""
    if piece.Type == "pawn" && (to.Row == 0 || to.Row == 7) {
        // Giả sử luôn phong thành hậu
        promotion = "=Q"
    }

    // Ăn quân
    captured := ce.GetPieceAt(gameBefore, to)
    isCapture := captured != nil
    // Ăn en passant
    if piece.Type == "pawn" && stateBefore.EnPassantSquare != nil &&
        to.Row == stateBefore.EnPassantSquare.Row && to.Col == stateBefore.EnPassantSquare.Col {
        isCapture = true
    }

    // Tên cột from (dùng cho tốt ăn)
    fromFile := string(rune('a' + from.Col))

    // Xây dựng notation
    notation := ""
    if piece.Type == "pawn" {
        if isCapture {
            notation += fromFile + "x" + ce.PositionToAlgebraic(to)
        } else {
            notation += ce.PositionToAlgebraic(to)
        }
    } else {
        notation += pieceSymbol
        // Nếu có nhiều quân cùng loại có thể đi đến ô đó, cần thêm disambiguation (bạn có thể bổ sung sau)
        if isCapture {
            notation += "x"
        }
        notation += ce.PositionToAlgebraic(to)
    }

    notation += promotion

    // Kiểm tra chiếu/chiếu hết
    opponentColor := "black"
    if piece.Color == "black" {
        opponentColor = "white"
    }
    if ce.IsInCheck(gameAfter, opponentColor) {
        if ce.IsCheckmate(gameAfter, stateAfter) {
            notation += "#"
        } else {
            notation += "+"
        }
    }

    return notation
}

func (ce *ChessEngine) ExecuteServerMove(state *ServerGameState, from Position, to Position) bool {
    // Validate move
    gameState := GameState{
    ActiveColor:     state.ActiveColor,
    CastlingRights:  state.CastlingRights,
    EnPassantSquare: state.EnPassantSquare,
    MoveCount:       state.FullMoveNumber,
    HalfMoveClock:   state.HalfMoveClock,
    }
    
    if !ce.ValidateMove(state.Bitboards, gameState, from, to, state.ActiveColor) {
        return false
    }
    
    // Get piece and check for capture
    piece := ce.GetPieceAt(state.Bitboards, from)
    capturedPiece := ce.GetPieceAt(state.Bitboards, to)

    // Update half-move clock
    if piece.Type == "pawn" || capturedPiece != nil {
        state.HalfMoveClock = 0
    } else {
        state.HalfMoveClock++
    }
    
    // Execute the move
    ce.ExecuteMove(&state.Bitboards, &gameState, from, to)
    
    // Update server state from game state
    state.ActiveColor = gameState.ActiveColor
    state.CastlingRights = gameState.CastlingRights
    state.EnPassantSquare = gameState.EnPassantSquare

    // Update full move number
    if state.ActiveColor == "white" {
        state.FullMoveNumber++
    }
    
    // Update material count
    ce.UpdateMaterialCount(state, capturedPiece)
    
    // Add move to history
    // moveNotation := ce.MoveToAlgebraic(state.Bitboards, from, to)
    // ce.AddMoveToHistory(state, moveNotation)
    
    // Update FEN and position counts
    newFen := ce.GameStateToFEN(state.Bitboards, gameState)
    ce.UpdatePositionCounts(state, newFen)
    
    return true
}

// Convert game state to FEN notation (simplified)
func (ce *ChessEngine) GameStateToFEN(game BitboardGame, state GameState) string {
    // This is a simplified version - you may want to implement full FEN conversion
    board := ce.BitboardsToFENBoard(game)
    activeColor := "w"
    if state.ActiveColor == "black" {
        activeColor = "b"
    }
    
    castling := ce.CastlingRightsToFEN(state.CastlingRights)
    enPassant := "-"
    if state.EnPassantSquare != nil {
        enPassant = ce.PositionToAlgebraic(*state.EnPassantSquare)
    }
    
    return fmt.Sprintf("%s %s %s %s 0 %d", board, activeColor, castling, enPassant, state.MoveCount)
}

// Helper: Convert bitboards to FEN board representation
func (ce *ChessEngine) BitboardsToFENBoard(game BitboardGame) string {
    // This is a placeholder - implement full FEN board conversion
    return "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
}

// Helper: Convert castling rights to FEN format
func (ce *ChessEngine) CastlingRightsToFEN(rights CastlingRights) string {
    result := ""
    if rights.WhiteKingSide {
        result += "K"
    }
    if rights.WhiteQueenSide {
        result += "Q"
    }
    if rights.BlackKingSide {
        result += "k"
    }
    if rights.BlackQueenSide {
        result += "q"
    }
    
    if result == "" {
        return "-"
    }
    return result
}

// Convert position to algebraic notation (for logging/debugging)
func (ce *ChessEngine) PositionToAlgebraic(pos Position) string {
    file := string(rune('a' + pos.Col))
    rank := fmt.Sprintf("%d", pos.Row+1)
    return file + rank
}

// Parse algebraic notation to position
func (ce *ChessEngine) AlgebraicToPosition(algebraic string) Position {
    if len(algebraic) != 2 {
        return Position{Row: -1, Col: -1}
    }
    
    col := int(algebraic[0] - 'a')
    row := 8 - int(algebraic[1] - '0')
    
    return Position{Row: row, Col: col}
}
