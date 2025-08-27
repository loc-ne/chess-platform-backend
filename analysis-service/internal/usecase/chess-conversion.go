package usecase

import (
    "github.com/notnil/chess"
    "strings"
    "regexp"
    "fmt"
)

func PGNtoUCI(moves []string) ([]string, error) {
    game := chess.NewGame()
    var uciMoves []string

    for _, m := range moves {
        move, err := findValidMove(game, m)
        if err != nil {
            return nil, err
        }
        
        uciMove := formatUCIMove(move)
        uciMoves = append(uciMoves, uciMove)
        game.Move(move)
    }

    return uciMoves, nil
}

func findValidMove(game *chess.Game, moveStr string) (*chess.Move, error) {
    validMoves := game.ValidMoves()
    
    // Try exact match first
    for _, mv := range validMoves {
        if mv.String() == moveStr {
            return mv, nil
        }
    }
    
    // Try case-insensitive match
    for _, mv := range validMoves {
        if strings.EqualFold(mv.String(), moveStr) {
            return mv, nil
        }
    }
    
    // Try algebraic notation
    notation := &chess.AlgebraicNotation{}
    if move, err := notation.Decode(game.Position(), moveStr); err == nil {
        return move, nil
    }
    
    return nil, fmt.Errorf("invalid move: %s", moveStr)
}

func formatUCIMove(move *chess.Move) string {
    from := move.S1().String()
    to := move.S2().String()
    
    promotion := ""
    if move.Promo() != chess.NoPieceType {
        switch move.Promo() {
        case chess.Queen:
            promotion = "q"
        case chess.Rook:
            promotion = "r"
        case chess.Bishop:
            promotion = "b"
        case chess.Knight:
            promotion = "n"
        }
    }
    
    return from + to + promotion
}

func ParsePGNToMoves(pgn string) []string {
    tokens := strings.Fields(pgn)
    var moves []string
    
    for i := 0; i < len(tokens); i++ {
        token := tokens[i]
        if matched, _ := regexp.MatchString(`\d+\.`, token); matched {
            if i+1 < len(tokens) {
                moves = append(moves, tokens[i+1])
                i++ 
            }
        } else {
            moves = append(moves, token)
        }
    }
    return moves
}