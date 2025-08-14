package engine

// Check if king is in check
func (ce *ChessEngine) IsInCheck(game BitboardGame, activeColor string) bool {
    var kingBitboard uint64
    if activeColor == "white" {
        kingBitboard = game.WhiteKing
    } else {
        kingBitboard = game.BlackKing
    }
    
    kingPositions := ce.ConvertBitboardToCoordinates(kingBitboard)
    if len(kingPositions) == 0 {
        return false
    }
    
    kingPos := kingPositions[0]
    enemyColor := "black"
    if activeColor == "black" {
        enemyColor = "white"
    }
    
    return ce.IsSquareAttackedBy(game, kingPos, enemyColor)
}

// Check if square is attacked by enemy
func (ce *ChessEngine) IsSquareAttackedBy(game BitboardGame, square Position, attackerColor string) bool {
    if attackerColor == "black" {
        if ce.IsPawnAttack(game.BlackPawns, square, "black") { return true }
        if ce.IsKnightAttack(game.BlackKnights, square) { return true }
        if ce.IsBishopAttack(game.BlackBishops, game, square) { return true }
        if ce.IsRookAttack(game.BlackRooks, game, square) { return true }
        if ce.IsQueenAttack(game.BlackQueens, game, square) { return true }
        if ce.IsKingAttack(game.BlackKing, square) { return true }
    } else {
        if ce.IsPawnAttack(game.WhitePawns, square, "white") { return true }
        if ce.IsKnightAttack(game.WhiteKnights, square) { return true }
        if ce.IsBishopAttack(game.WhiteBishops, game, square) { return true }
        if ce.IsRookAttack(game.WhiteRooks, game, square) { return true }
        if ce.IsQueenAttack(game.WhiteQueens, game, square) { return true }
        if ce.IsKingAttack(game.WhiteKing, square) { return true }
    }
    
    return false
}

// Check pawn attacks
func (ce *ChessEngine) IsPawnAttack(pawns uint64, targetSquare Position, attackerColor string) bool {
    var pawnAttackMoves []Position
    
    if attackerColor == "black" {
        pawnAttackMoves = []Position{
            {Row: targetSquare.Row - 1, Col: targetSquare.Col - 1},
            {Row: targetSquare.Row - 1, Col: targetSquare.Col + 1},
        }
    } else {
        pawnAttackMoves = []Position{
            {Row: targetSquare.Row + 1, Col: targetSquare.Col - 1},
            {Row: targetSquare.Row + 1, Col: targetSquare.Col + 1},
        }
    }
    
    for _, move := range pawnAttackMoves {
        if ce.IsValidSquare(move) {
            attackPos := move.Row*8 + move.Col
            attackBit := uint64(1) << uint(attackPos)
            if (pawns & attackBit) != 0 {
                return true
            }
        }
    }
    
    return false
}

// Check knight attacks
func (ce *ChessEngine) IsKnightAttack(knights uint64, targetSquare Position) bool {
    knightMoves := []Position{
        {Row: targetSquare.Row - 2, Col: targetSquare.Col - 1},
        {Row: targetSquare.Row - 2, Col: targetSquare.Col + 1},
        {Row: targetSquare.Row - 1, Col: targetSquare.Col - 2},
        {Row: targetSquare.Row - 1, Col: targetSquare.Col + 2},
        {Row: targetSquare.Row + 1, Col: targetSquare.Col - 2},
        {Row: targetSquare.Row + 1, Col: targetSquare.Col + 2},
        {Row: targetSquare.Row + 2, Col: targetSquare.Col - 1},
        {Row: targetSquare.Row + 2, Col: targetSquare.Col + 1},
    }
    
    for _, move := range knightMoves {
        if ce.IsValidSquare(move) {
            attackPos := move.Row*8 + move.Col
            attackBit := uint64(1) << uint(attackPos)
            if (knights & attackBit) != 0 {
                return true
            }
        }
    }
    
    return false
}

// Check bishop attacks
func (ce *ChessEngine) IsBishopAttack(bishops uint64, game BitboardGame, targetSquare Position) bool {
    allPieces := ce.GetAllPieces(game)
    diagonalDirections := []Position{
        {Row: -1, Col: -1}, {Row: -1, Col: 1},
        {Row: 1, Col: -1}, {Row: 1, Col: 1},
    }
    
    for _, direction := range diagonalDirections {
        currentRow := targetSquare.Row + direction.Row
        currentCol := targetSquare.Col + direction.Col
        
        for ce.IsValidSquare(Position{Row: currentRow, Col: currentCol}) {
            currentPos := currentRow*8 + currentCol
            currentBit := uint64(1) << uint(currentPos)
            
            if (allPieces & currentBit) != 0 {
                if (bishops & currentBit) != 0 {
                    return true
                }
                break
            }
            
            currentRow += direction.Row
            currentCol += direction.Col
        }
    }
    
    return false
}

// Check rook attacks
func (ce *ChessEngine) IsRookAttack(rooks uint64, game BitboardGame, targetSquare Position) bool {
    allPieces := ce.GetAllPieces(game)
    straightDirections := []Position{
        {Row: -1, Col: 0}, {Row: 1, Col: 0},
        {Row: 0, Col: -1}, {Row: 0, Col: 1},
    }
    
    for _, direction := range straightDirections {
        currentRow := targetSquare.Row + direction.Row
        currentCol := targetSquare.Col + direction.Col
        
        for ce.IsValidSquare(Position{Row: currentRow, Col: currentCol}) {
            currentPos := currentRow*8 + currentCol
            currentBit := uint64(1) << uint(currentPos)
            
            if (allPieces & currentBit) != 0 {
                if (rooks & currentBit) != 0 {
                    return true
                }
                break
            }
            
            currentRow += direction.Row
            currentCol += direction.Col
        }
    }
    
    return false
}

// Check queen attacks (rook + bishop)
func (ce *ChessEngine) IsQueenAttack(queens uint64, game BitboardGame, targetSquare Position) bool {
    return ce.IsRookAttack(queens, game, targetSquare) || ce.IsBishopAttack(queens, game, targetSquare)
}

// Check king attacks
func (ce *ChessEngine) IsKingAttack(king uint64, targetSquare Position) bool {
    kingMoves := []Position{
        {Row: targetSquare.Row - 1, Col: targetSquare.Col - 1},
        {Row: targetSquare.Row - 1, Col: targetSquare.Col},
        {Row: targetSquare.Row - 1, Col: targetSquare.Col + 1},
        {Row: targetSquare.Row, Col: targetSquare.Col - 1},
        {Row: targetSquare.Row, Col: targetSquare.Col + 1},
        {Row: targetSquare.Row + 1, Col: targetSquare.Col - 1},
        {Row: targetSquare.Row + 1, Col: targetSquare.Col},
        {Row: targetSquare.Row + 1, Col: targetSquare.Col + 1},
    }
    
    for _, move := range kingMoves {
        if ce.IsValidSquare(move) {
            pos := move.Row*8 + move.Col
            bit := uint64(1) << uint(pos)
            if (king & bit) != 0 {
                return true
            }
        }
    }
    
    return false
}
