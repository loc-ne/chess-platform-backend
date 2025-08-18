package engine

//  import ("fmt")
// Make move (basic implementation)
func (ce *ChessEngine) MakeMove(game *BitboardGame, from Position, to Position) {
    piece := ce.GetPieceAt(*game, from)
    if piece == nil {
        return
    }
    
    // Clear piece from source
    ce.ClearPieceAt(game, from, *piece)
    
    // Handle capture
    capturedPiece := ce.GetPieceAt(*game, to)
    if capturedPiece != nil {
        ce.ClearPieceAt(game, to, *capturedPiece)
    }
    
    // Place piece at destination
    ce.SetPieceAt(game, to, *piece)
}

// Update game state after move
func (ce *ChessEngine) UpdateGameState(game *BitboardGame, state *GameState, from Position, to Position, piece *Piece) {
    // Switch active color
    if state.ActiveColor == "black" {
        state.MoveCount++
    }
    
    // Update castling rights
    ce.UpdateCastlingRights(state, from, to, *piece)
    
    // Update en passant square
    ce.UpdateEnPassantSquare(state, from, to, *piece)
}

// Update castling rights based on move
func (ce *ChessEngine) UpdateCastlingRights(state *GameState, from Position, to Position, piece Piece) {
    // King moves - lose all castling rights for that color
    if piece.Type == "king" {
        if piece.Color == "white" {
            state.CastlingRights.WhiteKingSide = false
            state.CastlingRights.WhiteQueenSide = false
        } else {
            state.CastlingRights.BlackKingSide = false
            state.CastlingRights.BlackQueenSide = false
        }
    }
    
    // Rook moves or rook captured - lose castling rights for that side
    if piece.Type == "rook" || (from.Row == 0 || from.Row == 7) {
        // White rooks
        if from.Row == 7 && from.Col == 0 {
            state.CastlingRights.WhiteQueenSide = false
        }
        if from.Row == 7 && from.Col == 7 {
            state.CastlingRights.WhiteKingSide = false
        }
        // Black rooks
        if from.Row == 0 && from.Col == 0 {
            state.CastlingRights.BlackQueenSide = false
        }
        if from.Row == 0 && from.Col == 7 {
            state.CastlingRights.BlackKingSide = false
        }
    }
    
    // Rook captured - check destination square
    if to.Row == 0 && to.Col == 0 {
        state.CastlingRights.BlackQueenSide = false
    }
    if to.Row == 0 && to.Col == 7 {
        state.CastlingRights.BlackKingSide = false
    }
    if to.Row == 7 && to.Col == 0 {
        state.CastlingRights.WhiteQueenSide = false
    }
    if to.Row == 7 && to.Col == 7 {
        state.CastlingRights.WhiteKingSide = false
    }
}

// Update en passant square
func (ce *ChessEngine) UpdateEnPassantSquare(state *GameState, from Position, to Position, piece Piece) {
    state.EnPassantSquare = nil
    
    // Pawn moves two squares forward
    if piece.Type == "pawn" && abs(to.Row - from.Row) == 2 {
        state.EnPassantSquare = &Position{
            Row: (from.Row + to.Row) / 2,
            Col: from.Col,
        }
    } 
}

// Complete move execution with state update
func (ce *ChessEngine) ExecuteMove(game *BitboardGame, state *GameState, from Position, to Position) bool {
    // Handle special moves
    piece := ce.GetPieceAt(*game, from)
    if piece != nil {
        // Handle castling
        if piece.Type == "king" && abs(to.Col - from.Col) == 2 {
            ce.ExecuteCastling(game, from, to)
        } else if piece.Type == "pawn" && state.EnPassantSquare != nil && 
                 to.Row == state.EnPassantSquare.Row && to.Col == state.EnPassantSquare.Col {
            ce.ExecuteEnPassant(game, from, to, piece.Color)
        } else {
            // Regular move
            ce.MakeMove(game, from, to)
        }
    }
    // Update game state
    ce.UpdateGameState(game, state, from, to, piece)
    
    return true
}

// Validate move (main server function)
func (ce *ChessEngine) ValidateMove(game BitboardGame, state GameState, from Position, to Position, playerColor string) bool {
    // Check if it's player's turn
    if state.ActiveColor != playerColor {
        return false
    }
    
    // Check if piece exists and belongs to player
    piece := ce.GetPieceAt(game, from)
    if piece == nil || piece.Color != playerColor {
        return false
    }
    
    // Generate valid moves for the piece
    validMoves := ce.GenerateMovesForPiece(game, state, from)
    // Check if target move is in valid moves
    for _, move := range validMoves {
        if move.Row == to.Row && move.Col == to.Col {
            return true
        }
    }
    
    return false
}
