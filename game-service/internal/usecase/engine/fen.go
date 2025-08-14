package engine

import (
    "fmt"
    "strings"
)

func (ce *ChessEngine) BitboardToFEN(game BitboardGame, activeColor, castling, enPassant string, halfmove, fullmove int) string {
    pieceMap := map[string]rune{
        "white_pawn":   'P',
        "white_knight": 'N',
        "white_bishop": 'B',
        "white_rook":   'R',
        "white_queen":  'Q',
        "white_king":   'K',
        "black_pawn":   'p',
        "black_knight": 'n',
        "black_bishop": 'b',
        "black_rook":   'r',
        "black_queen":  'q',
        "black_king":   'k',
    }

    // Tạo mảng 8x8 lưu ký hiệu quân cờ
    board := [8][8]rune{}
    for row := 0; row < 8; row++ {
        for col := 0; col < 8; col++ {
            pos := Position{Row: row, Col: col}
            piece := ce.GetPieceAt(game, pos)
            if piece != nil {
                key := piece.Color + "_" + piece.Type
                board[row][col] = pieceMap[key]
            } else {
                board[row][col] = '.'
            }
        }
    }

    // Chuyển sang FEN
    fenRows := make([]string, 8)
    for row := 0; row < 8; row++ {
        fenRow := ""
        empty := 0
        for col := 0; col < 8; col++ {
            if board[7-row][col] == '.' { // <-- dùng board[7-row] thay vì board[row]
                empty++
            } else {
                if empty > 0 {
                    fenRow += fmt.Sprintf("%d", empty)
                    empty = 0
                }
                fenRow += string(board[7-row][col])
            }
        }
        if empty > 0 {
            fenRow += fmt.Sprintf("%d", empty)
        }
        fenRows[row] = fenRow
    }


    return fmt.Sprintf("%s %s %s %s %d %d",
        strings.Join(fenRows, "/"),
        string(activeColor[0]), castling, enPassant, halfmove, fullmove)
}

func (ce *ChessEngine) FENToBitboard(fen string) BitboardGame {
    var game BitboardGame
    pieceMap := map[rune]Piece{
        'P': {Type: "pawn", Color: "white"},
        'N': {Type: "knight", Color: "white"},
        'B': {Type: "bishop", Color: "white"},
        'R': {Type: "rook", Color: "white"},
        'Q': {Type: "queen", Color: "white"},
        'K': {Type: "king", Color: "white"},
        'p': {Type: "pawn", Color: "black"},
        'n': {Type: "knight", Color: "black"},
        'b': {Type: "bishop", Color: "black"},
        'r': {Type: "rook", Color: "black"},
        'q': {Type: "queen", Color: "black"},
        'k': {Type: "king", Color: "black"},
    }

    parts := strings.Split(fen, " ")
    rows := strings.Split(parts[0], "/")
    for row := 0; row < 8; row++ {
        col := 0
        for _, ch := range rows[row] {
            if ch >= '1' && ch <= '8' {
                col += int(ch - '0')
            } else {
                pos := Position{Row: row, Col: col}
                piece := pieceMap[ch]
                ce.SetPieceAt(&game, pos, piece)
                col++
            }
        }
    }
    return game
}

func (c CastlingRights) ToFEN() string {
    s := ""
    if c.WhiteKingSide { s += "K" }
    if c.WhiteQueenSide { s += "Q" }
    if c.BlackKingSide { s += "k" }
    if c.BlackQueenSide { s += "q" }
    if s == "" { s = "-" }
    return s
}

func (p *Position) ToFEN() string {
    if p == nil {
        return "-"
    }
    file := string('a' + p.Col)
    rank := fmt.Sprintf("%d", 8 - p.Row)
    return file + rank
}
