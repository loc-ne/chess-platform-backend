package engine


// Create starting position
func (ce *ChessEngine) CreateBitboardGame() BitboardGame {
    return BitboardGame{
        WhitePawns:   0x000000000000FF00,
        WhiteRooks:   0x0000000000000081,
        WhiteKnights: 0x0000000000000042,
        WhiteBishops: 0x0000000000000024,
        WhiteQueens:  0x0000000000000008,
        WhiteKing:    0x0000000000000010,
        
        BlackPawns:   0x00FF000000000000,
        BlackRooks:   0x8100000000000000,
        BlackKnights: 0x4200000000000000,
        BlackBishops: 0x2400000000000000,
        BlackQueens:  0x0800000000000000,
        BlackKing:    0x1000000000000000,
    }
}

// Helper: Get LSB position
func (ce *ChessEngine) GetLSBPosition(bitboard uint64) int {
    if bitboard == 0 {
        return -1
    }
    
    position := 0
    temp := bitboard
    
    for (temp & 1) == 0 {
        temp >>= 1
        position++
    }
    
    return position
}

// Helper: Clear LSB
func (ce *ChessEngine) ClearLSB(bitboard uint64) uint64 {
    return bitboard & (bitboard - 1)
}

// Convert bitboard to coordinates
func (ce *ChessEngine) ConvertBitboardToCoordinates(bitboard uint64) []Position {
    coordinates := []Position{}
    temp := bitboard
    
    for temp != 0 {
        position := ce.GetLSBPosition(temp)
        row := position / 8
        col := position % 8
        
        coordinates = append(coordinates, Position{Row: row, Col: col})
        temp = ce.ClearLSB(temp)
    }
    
    return coordinates
}

// Get all pieces bitboard
func (ce *ChessEngine) GetAllPieces(game BitboardGame) uint64 {
    return game.WhitePawns | game.WhiteRooks | game.WhiteKnights |
           game.WhiteBishops | game.WhiteQueens | game.WhiteKing |
           game.BlackPawns | game.BlackRooks | game.BlackKnights |
           game.BlackBishops | game.BlackQueens | game.BlackKing
}

// Get all pieces of specific color
func (ce *ChessEngine) GetAllPiecesOfColor(game BitboardGame, color string) uint64 {
    if color == "white" {
        return game.WhitePawns | game.WhiteRooks | game.WhiteKnights |
               game.WhiteBishops | game.WhiteQueens | game.WhiteKing
    }
    return game.BlackPawns | game.BlackRooks | game.BlackKnights |
           game.BlackBishops | game.BlackQueens | game.BlackKing
}

// Get piece at position
func (ce *ChessEngine) GetPieceAt(game BitboardGame, pos Position) *Piece {
    bit := uint64(1) << uint(pos.Row*8 + pos.Col)
    
    if (game.WhitePawns & bit) != 0 { return &Piece{Type: "pawn", Color: "white"} }
    if (game.BlackPawns & bit) != 0 { return &Piece{Type: "pawn", Color: "black"} }
    if (game.WhiteKnights & bit) != 0 { return &Piece{Type: "knight", Color: "white"} }
    if (game.BlackKnights & bit) != 0 { return &Piece{Type: "knight", Color: "black"} }
    if (game.WhiteBishops & bit) != 0 { return &Piece{Type: "bishop", Color: "white"} }
    if (game.BlackBishops & bit) != 0 { return &Piece{Type: "bishop", Color: "black"} }
    if (game.WhiteRooks & bit) != 0 { return &Piece{Type: "rook", Color: "white"} }
    if (game.BlackRooks & bit) != 0 { return &Piece{Type: "rook", Color: "black"} }
    if (game.WhiteQueens & bit) != 0 { return &Piece{Type: "queen", Color: "white"} }
    if (game.BlackQueens & bit) != 0 { return &Piece{Type: "queen", Color: "black"} }
    if (game.WhiteKing & bit) != 0 { return &Piece{Type: "king", Color: "white"} }
    if (game.BlackKing & bit) != 0 { return &Piece{Type: "king", Color: "black"} }
    
    return nil
}

// Check if position is valid
func (ce *ChessEngine) IsValidSquare(pos Position) bool {
    return pos.Row >= 0 && pos.Row < 8 && pos.Col >= 0 && pos.Col < 8
}

// Check if square is occupied by pieces bitboard
func (ce *ChessEngine) IsSquareOccupied(pieces uint64, pos Position) bool {
    bit := uint64(1) << uint(pos.Row*8 + pos.Col)
    return (pieces & bit) != 0
}

// Clone bitboard game
func (ce *ChessEngine) CloneBitboards(game BitboardGame) BitboardGame {
    return BitboardGame{
        WhitePawns:   game.WhitePawns,
        WhiteRooks:   game.WhiteRooks,
        WhiteKnights: game.WhiteKnights,
        WhiteBishops: game.WhiteBishops,
        WhiteQueens:  game.WhiteQueens,
        WhiteKing:    game.WhiteKing,
        BlackPawns:   game.BlackPawns,
        BlackRooks:   game.BlackRooks,
        BlackKnights: game.BlackKnights,
        BlackBishops: game.BlackBishops,
        BlackQueens:  game.BlackQueens,
        BlackKing:    game.BlackKing,
    }
}

// Clear piece at position
func (ce *ChessEngine) ClearPieceAt(game *BitboardGame, pos Position, piece Piece) {
    bit := uint64(1) << uint(pos.Row*8 + pos.Col)
    clearBit := ^bit
    
    if piece.Color == "white" {
        switch piece.Type {
        case "pawn": game.WhitePawns &= clearBit
        case "knight": game.WhiteKnights &= clearBit
        case "bishop": game.WhiteBishops &= clearBit
        case "rook": game.WhiteRooks &= clearBit
        case "queen": game.WhiteQueens &= clearBit
        case "king": game.WhiteKing &= clearBit
        }
    } else {
        switch piece.Type {
        case "pawn": game.BlackPawns &= clearBit
        case "knight": game.BlackKnights &= clearBit
        case "bishop": game.BlackBishops &= clearBit
        case "rook": game.BlackRooks &= clearBit
        case "queen": game.BlackQueens &= clearBit
        case "king": game.BlackKing &= clearBit
        }
    }
}

// Set piece at position
func (ce *ChessEngine) SetPieceAt(game *BitboardGame, pos Position, piece Piece) {
    bit := uint64(1) << uint(pos.Row*8 + pos.Col)
    
    if piece.Color == "white" {
        switch piece.Type {
        case "pawn": game.WhitePawns |= bit
        case "knight": game.WhiteKnights |= bit
        case "bishop": game.WhiteBishops |= bit
        case "rook": game.WhiteRooks |= bit
        case "queen": game.WhiteQueens |= bit
        case "king": game.WhiteKing |= bit
        }
    } else {
        switch piece.Type {
        case "pawn": game.BlackPawns |= bit
        case "knight": game.BlackKnights |= bit
        case "bishop": game.BlackBishops |= bit
        case "rook": game.BlackRooks |= bit
        case "queen": game.BlackQueens |= bit
        case "king": game.BlackKing |= bit
        }
    }
}

// Count set bits in bitboard
func (ce *ChessEngine) CountBits(bitboard uint64) int {
    count := 0
    for bitboard != 0 {
        count++
        bitboard &= bitboard - 1 // Clear least significant bit
    }
    return count
}

// Helper: absolute value
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
