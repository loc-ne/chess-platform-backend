package engine

// import ("fmt")

// Castling helper methods
func (ce *ChessEngine) CanCastleKingside(game BitboardGame, state GameState, color string) bool {
    isWhite := color == "white"

    // Check castling rights
    if isWhite && !state.CastlingRights.WhiteKingSide {
        return false
    }
    if !isWhite && !state.CastlingRights.BlackKingSide {
        return false
    }
    
    kingRow := 0
    if !isWhite {
        kingRow = 7
    }
    
    // Check if squares between king and rook are empty
    squaresBetween := []Position{
        {Row: kingRow, Col: 5}, // f1 or f8
        {Row: kingRow, Col: 6}, // g1 or g8
    }
    
    allPieces := ce.GetAllPieces(game)
    enemyColor := "black"
    if !isWhite {
        enemyColor = "white"
    }

    for _, square := range squaresBetween {
        // Square must be empty
        if ce.IsSquareOccupied(allPieces, square) {
            return false
        }
        
        // Square must not be attacked by enemy
        if ce.IsSquareAttackedBy(game, square, enemyColor) {
            return false
        }
    }
    
    // Check if rook is in correct position
    rookPosition := Position{Row: kingRow, Col: 7}
    rookPiece := ce.GetPieceAt(game, rookPosition)
    
    return rookPiece != nil && rookPiece.Type == "rook" && rookPiece.Color == color
}

func (ce *ChessEngine) CanCastleQueenside(game BitboardGame, state GameState, color string) bool {
    isWhite := color == "white"
    
    // Check castling rights
    if isWhite && !state.CastlingRights.WhiteQueenSide {
        return false
    }
    if !isWhite && !state.CastlingRights.BlackQueenSide {
        return false
    }
    
    kingRow := 0
    if !isWhite {
        kingRow = 7
    }
    
    // Check if squares between king and rook are empty
    squaresBetween := []Position{
        {Row: kingRow, Col: 1}, // b1 or b8
        {Row: kingRow, Col: 2}, // c1 or c8
        {Row: kingRow, Col: 3}, // d1 or d8
    }
    
    // Squares king moves through (must not be attacked)
    squaresKingMovesThrough := []Position{
        {Row: kingRow, Col: 2}, // c1 or c8
        {Row: kingRow, Col: 3}, // d1 or d8
    }
    
    allPieces := ce.GetAllPieces(game)
    enemyColor := "black"
    if !isWhite {
        enemyColor = "white"
    }
    
    // Check all squares between are empty
    for _, square := range squaresBetween {
        if ce.IsSquareOccupied(allPieces, square) {
            return false
        }
    }
    
    // Check squares king moves through are not attacked
    for _, square := range squaresKingMovesThrough {
        if ce.IsSquareAttackedBy(game, square, enemyColor) {
            return false
        }
    }
    
    // Check if rook is in correct position
    rookPosition := Position{Row: kingRow, Col: 0}
    rookPiece := ce.GetPieceAt(game, rookPosition)
    
    return rookPiece != nil && rookPiece.Type == "rook" && rookPiece.Color == color
}

// Execute castling move
func (ce *ChessEngine) ExecuteCastling(game *BitboardGame, kingFrom Position, kingTo Position) {
    // Move king
    ce.MakeMove(game, kingFrom, kingTo)
    
    // Move rook
    var rookFrom, rookTo Position
    if kingTo.Col == 6 { // Kingside castling
        rookFrom = Position{Row: kingFrom.Row, Col: 7}
        rookTo = Position{Row: kingFrom.Row, Col: 5}
    } else { 
        rookFrom = Position{Row: kingFrom.Row, Col: 0}
        rookTo = Position{Row: kingFrom.Row, Col: 3}
    }
    
    ce.MakeMove(game, rookFrom, rookTo)
}

// Execute en passant capture
func (ce *ChessEngine) ExecuteEnPassant(game *BitboardGame, from Position, to Position, color string) {
    // Move pawn
    ce.MakeMove(game, from, to)
    
    // Remove captured pawn
    capturedPawnRow := from.Row
    capturedPawnPos := Position{Row: capturedPawnRow, Col: to.Col}
    capturedPawn := Piece{Type: "pawn", Color: map[string]string{"white": "black", "black": "white"}[color]}
    ce.ClearPieceAt(game, capturedPawnPos, capturedPawn)
}
