package engine

// Generate moves for piece (simplified version for server validation)
func (ce *ChessEngine) GenerateMovesForPiece(game BitboardGame, state GameState, fromPos Position) []Position {
    piece := ce.GetPieceAt(game, fromPos)
    if piece == nil {
        return []Position{}
    }
    
    var moves []Position
    
    switch piece.Type {
    case "pawn":
        moves = ce.GeneratePawnMoves(game, state, fromPos)
    case "knight":
        moves = ce.GenerateKnightMoves(game, fromPos)
    case "bishop":
        moves = ce.GenerateBishopMoves(game, fromPos)
    case "rook":
        moves = ce.GenerateRookMoves(game, fromPos)
    case "queen":
        moves = ce.GenerateQueenMoves(game, fromPos)
    case "king":
        moves = ce.GenerateKingMoves(game, state, fromPos)
    }
    
    return ce.FilterLegalMoves(game, state, fromPos, moves)
}

// Filter legal moves (remove moves that leave king in check)
func (ce *ChessEngine) FilterLegalMoves(game BitboardGame, state GameState, fromPos Position, moves []Position) []Position {
    legalMoves := []Position{}
    
    for _, move := range moves {
        newGame := ce.CloneBitboards(game)
        ce.MakeMove(&newGame, fromPos, move)
        
        if !ce.IsInCheck(newGame, state.ActiveColor) {
            legalMoves = append(legalMoves, move)
        }
    }
    
    return legalMoves
}

// Generate knight moves
func (ce *ChessEngine) GenerateKnightMoves(game BitboardGame, from Position) []Position {
    moves := []Position{}
    piece := ce.GetPieceAt(game, from)
    if piece == nil || piece.Type != "knight" {
        return moves
    }
    
    ownPieces := ce.GetAllPiecesOfColor(game, piece.Color)
    knightMoves := []Position{
        {Row: from.Row - 2, Col: from.Col - 1}, {Row: from.Row - 2, Col: from.Col + 1},
        {Row: from.Row - 1, Col: from.Col - 2}, {Row: from.Row - 1, Col: from.Col + 2},
        {Row: from.Row + 1, Col: from.Col - 2}, {Row: from.Row + 1, Col: from.Col + 2},
        {Row: from.Row + 2, Col: from.Col - 1}, {Row: from.Row + 2, Col: from.Col + 1},
    }
    
    for _, move := range knightMoves {
        if ce.IsValidSquare(move) && !ce.IsSquareOccupied(ownPieces, move) {
            moves = append(moves, move)
        }
    }
    
    return moves
}

// Generate bishop moves
func (ce *ChessEngine) GenerateBishopMoves(game BitboardGame, from Position) []Position {
    validMoves := []Position{}
    piece := ce.GetPieceAt(game, from)
    
    if piece == nil || piece.Type != "bishop" {
        return validMoves
    }
    
    allPieces := ce.GetAllPieces(game)
    ownPieces := ce.GetAllPiecesOfColor(game, piece.Color)
    
    // 4 diagonal directions
    diagonalDirections := []Position{
        {Row: -1, Col: -1}, // Up-left
        {Row: -1, Col: 1},  // Up-right
        {Row: 1, Col: -1},  // Down-left
        {Row: 1, Col: 1},   // Down-right
    }
    
    for _, direction := range diagonalDirections {
        step := 1
        
        for {
            newRow := from.Row + step*direction.Row
            newCol := from.Col + step*direction.Col
            newPosition := Position{Row: newRow, Col: newCol}
            
            // Check bounds
            if !ce.IsValidSquare(newPosition) {
                break
            }
            
            bit := uint64(1) << uint(newRow*8 + newCol)
            
            // Check if square is occupied
            if (allPieces & bit) != 0 {
                // Found a piece
                if !ce.IsSquareOccupied(ownPieces, newPosition) {
                    // Enemy piece - can capture
                    validMoves = append(validMoves, newPosition)
                }
                // Stop ray (can't jump over pieces)
                break
            } else {
                // Empty square - can move here
                validMoves = append(validMoves, newPosition)
            }
            
            step++
        }
    }
    
    return validMoves
}

// Generate rook moves
func (ce *ChessEngine) GenerateRookMoves(game BitboardGame, from Position) []Position {
    validMoves := []Position{}
    piece := ce.GetPieceAt(game, from)
    
    if piece == nil || piece.Type != "rook" {
        return validMoves
    }
    
    allPieces := ce.GetAllPieces(game)
    ownPieces := ce.GetAllPiecesOfColor(game, piece.Color)
    
    // 4 straight directions
    straightDirections := []Position{
        {Row: -1, Col: 0}, // Up
        {Row: 1, Col: 0},  // Down
        {Row: 0, Col: -1}, // Left
        {Row: 0, Col: 1},  // Right
    }
    
    for _, direction := range straightDirections {
        step := 1
        
        for {
            newRow := from.Row + step*direction.Row
            newCol := from.Col + step*direction.Col
            newPosition := Position{Row: newRow, Col: newCol}
            
            // Check bounds
            if !ce.IsValidSquare(newPosition) {
                break
            }
            
            bit := uint64(1) << uint(newRow*8 + newCol)
            
            // Check if square is occupied
            if (allPieces & bit) != 0 {
                // Found a piece
                if !ce.IsSquareOccupied(ownPieces, newPosition) {
                    // Enemy piece - can capture
                    validMoves = append(validMoves, newPosition)
                }
                // Stop ray (can't jump over pieces)
                break
            } else {
                // Empty square - can move here
                validMoves = append(validMoves, newPosition)
            }
            
            step++
        }
    }
    
    return validMoves
}

// Generate queen moves (rook + bishop)
func (ce *ChessEngine) GenerateQueenMoves(game BitboardGame, from Position) []Position {
    validMoves := []Position{}
    piece := ce.GetPieceAt(game, from)
    
    if piece == nil || piece.Type != "queen" {
        return validMoves
    }
    
    allPieces := ce.GetAllPieces(game)
    ownPieces := ce.GetAllPiecesOfColor(game, piece.Color)
    
    // Queen = Rook + Bishop directions (8 total)
    allDirections := []Position{
        // Rook directions (straight)
        {Row: -1, Col: 0}, // Up
        {Row: 1, Col: 0},  // Down
        {Row: 0, Col: -1}, // Left
        {Row: 0, Col: 1},  // Right
        
        // Bishop directions (diagonal)
        {Row: -1, Col: -1}, // Up-left
        {Row: -1, Col: 1},  // Up-right
        {Row: 1, Col: -1},  // Down-left
        {Row: 1, Col: 1},   // Down-right
    }
    
    for _, direction := range allDirections {
        step := 1
        
        for {
            newRow := from.Row + step*direction.Row
            newCol := from.Col + step*direction.Col
            newPosition := Position{Row: newRow, Col: newCol}
            
            // Check bounds
            if !ce.IsValidSquare(newPosition) {
                break
            }
            
            bit := uint64(1) << uint(newRow*8 + newCol)
            
            // Check if square is occupied
            if (allPieces & bit) != 0 {
                // Found a piece
                if !ce.IsSquareOccupied(ownPieces, newPosition) {
                    // Enemy piece - can capture
                    validMoves = append(validMoves, newPosition)
                }
                // Stop ray (can't jump over pieces)
                break
            } else {
                // Empty square - can move here
                validMoves = append(validMoves, newPosition)
            }
            
            step++
        }
    }
    
    return validMoves
}

// Generate king moves
func (ce *ChessEngine) GenerateKingMoves(game BitboardGame, state GameState, from Position) []Position {
    validMoves := []Position{}
    piece := ce.GetPieceAt(game, from)
    
    if piece == nil || piece.Type != "king" {
        return validMoves
    }
    
    ownPieces := ce.GetAllPiecesOfColor(game, piece.Color)
    enemyColor := "black"
    if piece.Color == "black" {
        enemyColor = "white"
    }
    
    // 8 directions (king can move 1 square in any direction)
    kingMoves := []Position{
        {Row: -1, Col: -1}, {Row: -1, Col: 0}, {Row: -1, Col: 1},
        {Row: 0, Col: -1}, {Row: 0, Col: 1},
        {Row: 1, Col: -1}, {Row: 1, Col: 0}, {Row: 1, Col: 1},
    }
    
    // Regular king moves (1 square in any direction)
    for _, move := range kingMoves {
        newRow := from.Row + move.Row
        newCol := from.Col + move.Col
        newPosition := Position{Row: newRow, Col: newCol}
        
        // Check bounds
        if !ce.IsValidSquare(newPosition) {
            continue
        }
        
        // Can't move to square occupied by own pieces
        if ce.IsSquareOccupied(ownPieces, newPosition) {
            continue
        }
        
        // Important: Can't move to square attacked by enemy
        if ce.IsSquareAttackedBy(game, newPosition, enemyColor) {
            continue
        }
        
        validMoves = append(validMoves, newPosition)
    }
    
    // Castling moves (if not in check and castling rights available)
    isInCheck := ce.IsInCheck(game, state.ActiveColor)
    if !isInCheck {
        // Kingside castling
        if ce.CanCastleKingside(game, state, piece.Color) {
            kingsideCastlePosition := Position{
                Row: from.Row,
                Col: from.Col + 2,
            }
            validMoves = append(validMoves, kingsideCastlePosition)
        }
        
        // Queenside castling
        if ce.CanCastleQueenside(game, state, piece.Color) {
            queensideCastlePosition := Position{
                Row: from.Row,
                Col: from.Col - 2,
            }
            validMoves = append(validMoves, queensideCastlePosition)
        }
    }
    
    return validMoves
}

// Improved pawn moves with all rules
func (ce *ChessEngine) GeneratePawnMoves(game BitboardGame, state GameState, from Position) []Position {
    validMoves := []Position{}
    piece := ce.GetPieceAt(game, from)
    
    if piece == nil || piece.Type != "pawn" {
        return validMoves
    }
    
    isWhite := piece.Color == "white"
    direction := 1
    if !isWhite {
        direction = -1 // Black moves down (+1)
    }
    
    startingRow := 1
    if !isWhite {
        startingRow = 6
    }
    
    allPieces := ce.GetAllPieces(game)
    enemyPieces := ce.GetAllPiecesOfColor(game, map[bool]string{true: "black", false: "white"}[isWhite])
    
    // Forward move (1 square)
    oneForward := Position{Row: from.Row + direction, Col: from.Col}
    if ce.IsValidSquare(oneForward) && !ce.IsSquareOccupied(allPieces, oneForward) {
        validMoves = append(validMoves, oneForward)
        
        // Two squares forward from starting position
        if from.Row == startingRow {
            twoForward := Position{Row: from.Row + 2*direction, Col: from.Col}
            if ce.IsValidSquare(twoForward) && !ce.IsSquareOccupied(allPieces, twoForward) {
                validMoves = append(validMoves, twoForward)
            }
        }
    }
    
    // Diagonal captures (only if enemy piece present)
    captureLeft := Position{Row: from.Row + direction, Col: from.Col - 1}
    captureRight := Position{Row: from.Row + direction, Col: from.Col + 1}
    
    if ce.IsValidSquare(captureLeft) && ce.IsSquareOccupied(enemyPieces, captureLeft) {
        validMoves = append(validMoves, captureLeft)
    }
    
    if ce.IsValidSquare(captureRight) && ce.IsSquareOccupied(enemyPieces, captureRight) {
        validMoves = append(validMoves, captureRight)
    }
    
    // En passant
    if state.EnPassantSquare != nil {
        enPassantTargetRow := from.Row + direction
        if state.EnPassantSquare.Row == enPassantTargetRow &&
            abs(state.EnPassantSquare.Col - from.Col) == 1 {
            validMoves = append(validMoves, *state.EnPassantSquare)
        }
    }
    
    return validMoves
}
