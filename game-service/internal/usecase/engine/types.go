package engine

type Position struct {
    Row int
    Col int
}

type Piece struct {
    Type  string // "pawn", "knight", "bishop", "rook", "queen", "king"
    Color string // "white", "black"
}

// Material count for each color
type MaterialCount struct {
    Pawns   int `json:"pawns"`
    Knights int `json:"knights"`
    Bishops int `json:"bishops"`
    Rooks   int `json:"rooks"`
    Queens  int `json:"queens"`
}

type MoveNotation struct {
    MoveNumber int    `json:"moveNumber"`
    White      string `json:"white"`
    Black      string `json:"black"`
}
// Castling rights
type CastlingRights struct {
    WhiteKingSide  bool `json:"whiteKingSide"`
    WhiteQueenSide bool `json:"whiteQueenSide"`
    BlackKingSide  bool `json:"blackKingSide"`
    BlackQueenSide bool `json:"blackQueenSide"`
}

// Game state for internal logic
type GameState struct {
    ActiveColor     string
    CastlingRights  CastlingRights
    EnPassantSquare *Position
    MoveCount       int
	HalfMoveClock   int
}

// Complete server game state
type ServerGameState struct {
    CurrentFen      string                 `json:"currentFen"`
    Bitboards       BitboardGame          `json:"bitboards"`
    ActiveColor     string                `json:"activeColor"`
    CastlingRights  CastlingRights        `json:"castlingRights"`
    EnPassantSquare *Position             `json:"enPassantSquare"`
    MoveHistory     []MoveNotation        `json:"moveHistory"`
    FullMoveNumber  int                   `json:"fullMoveNumber"`
    HalfMoveClock   int                   `json:"halfMoveClock"`
    PositionCounts  map[string]int        `json:"positionCounts"`
    MaterialCount   map[string]MaterialCount `json:"materialCount"`
}

type ClientGameState struct {
    CurrentFen      string                 `json:"currentFen"`
    Bitboards       BitboardGame          `json:"bitboards"`
    ActiveColor     string                `json:"activeColor"`
    CastlingRights  CastlingRights        `json:"castlingRights"`
    EnPassantSquare *Position             `json:"enPassantSquare"`
}

// Piece count for insufficient material check
type PieceCount struct {
    Total   int
    Pawns   int
    Knights int
    Bishops int
    Rooks   int
    Queens  int
}

// Chess Engine
type ChessEngine struct{}
