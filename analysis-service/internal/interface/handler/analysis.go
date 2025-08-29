package handler

import (
	"fmt"          
    "strings"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/locne/analysis-service/internal/usecase"
)

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  []string    `json:"errors,omitempty"`
}

type AnalyzePGNRequest struct {
    PGN string `json:"pgn" binding:"required"`
}

type MoveLabel string

const (
    Brilliant  MoveLabel = "brilliant"
    Best      MoveLabel = "best"
    Excellent MoveLabel = "excellent"  // Nước đi xuất sắc
    Good      MoveLabel = "good"       // Nước đi tốt
    Inaccuracy MoveLabel = "inaccuracy" // Thiếu chính xác
    Mistake   MoveLabel = "mistake"     // Sai lầm
    Blunder   MoveLabel = "blunder"     // Sai lầm nghiêm trọng
    Book      MoveLabel = "book"        // Nước đi theo lý thuyết
    Miss      MoveLabel = "miss"
    Great     MoveLabel = "great"
)

type DetailedMoveAnalysis struct {
    MoveNumber int       `json:"move_number"`
    Player     string    `json:"player"`
    Move       string    `json:"move"`
    From       string    `json:"from"`
    To         string    `json:"to"`
    Label      MoveLabel `json:"label"`
    Evaluation float64   `json:"evaluation"`
}

type DetailedGameAnalysis struct {
    OpeningName string                 `json:"opening_name"`
    Moves       []DetailedMoveAnalysis `json:"moves"`
}

type PlayerAccuracy struct {
    Player       string  `json:"player"`
    Accuracy     float64 `json:"accuracy"`   
    BestMoveCount int   `json:"best_count"`  
    ExcellentCount int   `json:"excellent_count"`
    GoodMoveCount    int     `json:"good_count"`
    InaccuracyMoveCount int  `json:"inaccuracy_count"`
    MistakeMoveCount int     `json:"mistake_count"`
    BlunderMoveCount int     `json:"blunder_count"`
    BookMoveCount    int     `json:"book_count"`
    MissMoveCount    int     `json:"miss_count"`
    GreatMoveCount    int     `json:"great_count"`
    BrilliantMoveCount    int     `json:"brilliant_count"`
}

func ConvertToDetailedAnalysis(analysis *usecase.GameAnalysis, moves []string) (*DetailedGameAnalysis, error) {
    if analysis == nil {
        return nil, fmt.Errorf("analysis cannot be nil")
    }

    uciMoves, err := usecase.PGNtoUCI(moves)
    if err != nil {
        return nil, fmt.Errorf("failed to convert PGN to UCI: %w", err)
    }

    detailedMoves := make([]DetailedMoveAnalysis, len(analysis.Moves))
    
    for i, moveAnalysis := range analysis.Moves {
        player := "white"
        if i%2 == 1 {
            player = "black"
        }

        moveNumber := (i / 2) + 1

        from, to := parseUCIMove(uciMoves[i])

        label := convertToMoveLabel(moveAnalysis.Label)

        detailedMoves[i] = DetailedMoveAnalysis{
            MoveNumber: moveNumber,
            Player:     player,
            Move:       moveAnalysis.Move,
            From:       from,
            To:         to,
            Label:      label,
            Evaluation: moveAnalysis.CPScore,
        }
    }

    return &DetailedGameAnalysis{
        OpeningName: analysis.OpeningName,
        Moves:       detailedMoves,
    }, nil
}

func parseUCIMove(uciMove string) (string, string) {
    if len(uciMove) < 4 {
        return "??", "??"
    }
    
    from := uciMove[:2]
    to := uciMove[2:4]
    
    return from, to
}

func convertToMoveLabel(label string) MoveLabel {
    switch strings.ToLower(label) {
    case "brilliant":      
        return Brilliant
    case "best":
        return Best        
    case "excellent":
        return Excellent
    case "good":
        return Good
    case "inaccuracy":
        return Inaccuracy
    case "mistake":
        return Mistake
    case "blunder":
        return Blunder
    case "book":
        return Book
    case "great":
        return Great
    case "miss":
        return Miss
    default:
        return Book 
    }
}


func AnalyzeGame(c *gin.Context, trie *usecase.OpeningTrie) {
    var req AnalyzePGNRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, APIResponse{
            Status:  "error",
            Message: "Invalid request format",
            Errors:  []string{err.Error()},
        })
        return
    }

    inputMoves := usecase.ParsePGNToMoves(req.PGN)
    
    analysis, err := usecase.AnalyzeGame(inputMoves, trie)
    fmt.Println("what: ", analysis)
    if err != nil {
        c.JSON(http.StatusInternalServerError, APIResponse{
            Status:  "error",
            Message: "Failed to analyze game",
            Errors:  []string{err.Error()},
        })
        return
    }

    data, err := ConvertToDetailedAnalysis(analysis, inputMoves)
    if err != nil {
        c.JSON(http.StatusInternalServerError, APIResponse{
            Status:  "error",
            Message: "Failed to convert analysis",
            Errors:  []string{err.Error()},
        })
        return
    }

    c.JSON(http.StatusOK, APIResponse{
        Status:  "success",
        Message: "PGN analyzed successfully",
        Data:    data,
    })
}

func RegisterAnalysisRoutes(router *gin.Engine, trie *usecase.OpeningTrie) {
    api := router.Group("/api/v1/analysis")
    {
        api.POST("/pgn", func(c *gin.Context) {
            AnalyzeGame(c, trie)
        })
    }
}