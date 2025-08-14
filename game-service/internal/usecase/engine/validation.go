package engine

// Check for checkmate
func (ce *ChessEngine) IsCheckmate(game BitboardGame, state GameState) bool {
    if !ce.IsInCheck(game, state.ActiveColor) {
        return false
    }
    
    // Find all possible moves for current player
    allMoves := ce.GetAllLegalMoves(game, state, state.ActiveColor)
    return len(allMoves) == 0
}

// Check for stalemate
func (ce *ChessEngine) IsStalemate(game BitboardGame, state GameState) bool {
    if ce.IsInCheck(game, state.ActiveColor) {
        return false
    }
    
    // Find all possible moves for current player
    allMoves := ce.GetAllLegalMoves(game, state, state.ActiveColor)
    return len(allMoves) == 0
}

// Get all legal moves for a color
func (ce *ChessEngine) GetAllLegalMoves(game BitboardGame, state GameState, color string) []Position {
    allMoves := []Position{}
    
    // Get all pieces of the color
    pieces := ce.GetAllPiecesOfColor(game, color)
    positions := ce.ConvertBitboardToCoordinates(pieces)
    
    for _, pos := range positions {
        moves := ce.GenerateMovesForPiece(game, state, pos)
        allMoves = append(allMoves, moves...)
    }
    
    return allMoves
}

// Check if game is over
func (ce *ChessEngine) IsGameOver(game BitboardGame, state GameState, positionCounts map[string]int, currentFen string) (bool, string) {
    // Automatic draws (no claim needed)
    if ce.IsThreefoldRepetition(positionCounts, currentFen) {
        return true, "fivefold repetition"
    }

	if ce.IsFiftyMoveRule(state.HalfMoveClock) {
		return true, "fifty move rule"
	}

    // Checkmate/Stalemate
    if ce.IsCheckmate(game, state) {
        return true, "checkmate"
    }
    
    if ce.IsStalemate(game, state) {
        return true, "stalemate"
    }
    
    // Insufficient material
    if ce.IsInsufficientMaterial(game) {
        return true, "insufficient material"
    }
    
    return false, ""
}
// Check for insufficient material
func (ce *ChessEngine) IsInsufficientMaterial(game BitboardGame) bool {
    whitePieces := ce.CountPieces(game, "white")
    blackPieces := ce.CountPieces(game, "black")
    
    // King vs King
    if whitePieces.Total == 1 && blackPieces.Total == 1 {
        return true
    }
    
    // King and Bishop vs King
    if (whitePieces.Total == 2 && whitePieces.Bishops == 1 && blackPieces.Total == 1) ||
       (blackPieces.Total == 2 && blackPieces.Bishops == 1 && whitePieces.Total == 1) {
        return true
    }
    
    // King and Knight vs King
    if (whitePieces.Total == 2 && whitePieces.Knights == 1 && blackPieces.Total == 1) ||
       (blackPieces.Total == 2 && blackPieces.Knights == 1 && whitePieces.Total == 1) {
        return true
    }
    
    // King and Bishop vs King and Bishop (same colored squares)
    if whitePieces.Total == 2 && whitePieces.Bishops == 1 &&
       blackPieces.Total == 2 && blackPieces.Bishops == 1 {
        return ce.BishopsOnSameColorSquares(game)
    }
    
    // King and Knights vs King (multiple knights cannot force checkmate)
    if (whitePieces.Total > 1 && whitePieces.Knights == whitePieces.Total-1 && blackPieces.Total == 1) ||
       (blackPieces.Total > 1 && blackPieces.Knights == blackPieces.Total-1 && whitePieces.Total == 1) {
        return true
    }
    
    return false
}

// Count pieces for each color
func (ce *ChessEngine) CountPieces(game BitboardGame, color string) PieceCount {
    count := PieceCount{}
    
    if color == "white" {
        count.Pawns = ce.CountBits(game.WhitePawns)
        count.Knights = ce.CountBits(game.WhiteKnights)
        count.Bishops = ce.CountBits(game.WhiteBishops)
        count.Rooks = ce.CountBits(game.WhiteRooks)
        count.Queens = ce.CountBits(game.WhiteQueens)
    } else {
        count.Pawns = ce.CountBits(game.BlackPawns)
        count.Knights = ce.CountBits(game.BlackKnights)
        count.Bishops = ce.CountBits(game.BlackBishops)
        count.Rooks = ce.CountBits(game.BlackRooks)
        count.Queens = ce.CountBits(game.BlackQueens)
    }
    
    count.Total = count.Pawns + count.Knights + count.Bishops + count.Rooks + count.Queens + 1 // +1 for king
    return count
}

func (ce *ChessEngine) IsFiftyMoveRule(halfMoveClock int) bool {
    return halfMoveClock >= 100 // 50 full moves = 100 half-moves
}

func (ce *ChessEngine) IsThreefoldRepetition(positionCounts map[string]int, currentFen string) bool {
    count, exists := positionCounts[currentFen]
    return exists && count >= 3
}

func (ce *ChessEngine) BishopsOnSameColorSquares(game BitboardGame) bool {
    // Get bishop positions
    whitebishopPos := ce.ConvertBitboardToCoordinates(game.WhiteBishops)
    blackBishopPos := ce.ConvertBitboardToCoordinates(game.BlackBishops)
    
    if len(whitebishopPos) != 1 || len(blackBishopPos) != 1 {
        return false
    }
    
    // Check if bishops are on same colored squares
    whiteSquareColor := (whitebishopPos[0].Row + whitebishopPos[0].Col) % 2
    blackSquareColor := (blackBishopPos[0].Row + blackBishopPos[0].Col) % 2
    
    return whiteSquareColor == blackSquareColor
}