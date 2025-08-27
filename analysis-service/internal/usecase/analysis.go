package usecase

import (
    "math"
    "fmt"
)

type MoveAnalysis struct {
    Move      string
    Label     string
    CPScore   float64
    Position  int
}

type GameAnalysis struct {
    Moves         []MoveAnalysis
    OpeningName   string
    OpeningEnd    int
}

type MoveEvaluation struct {
    Score       float64
    IsMate      bool
    MateIn      int
    IsWhiteMove bool
}

func AnalyzeGame(inputMoves []string, trie *OpeningTrie) (*GameAnalysis, error) {
    engine, err := NewStockfishEngine()
    if err != nil {
        return nil, err
    }
    defer engine.Close()

    analysis := &GameAnalysis{
        Moves: make([]MoveAnalysis, len(inputMoves)),
    }

    // Analyze opening phase
    openingEnd := analyzeOpeningPhase(inputMoves, trie, analysis, engine)
    
    // Analyze game phase
    if openingEnd < len(inputMoves) {
        err := analyzeGamePhase(inputMoves, openingEnd, engine, analysis)
        if err != nil {
            return nil, err
        }
    }

    return analysis, nil
}

func analyzeOpeningPhase(inputMoves []string, trie *OpeningTrie, analysis *GameAnalysis, engine *StockfishEngine) int {
    openingEnd := 0
    var lastOpeningName string
    
    for i := 1; i <= len(inputMoves); i++ {
        currentMoves := inputMoves[:i]
        openings := trie.Search(currentMoves)
        engine.SetPosition(currentMoves)

        score, err, isMate, mateIn := engine.GetCurrentScore()
        if err != nil {
            score = 0
        }
            
        isWhiteMove := (i - 1) % 2 == 0
        moveEvaluation := createEvaluation(score, isMate, mateIn, isWhiteMove)

        if len(openings) > 0 {
            openingEnd = i
            lastOpeningName = openings[0].Name
            analysis.Moves[i - 1] = MoveAnalysis{
            Move:     inputMoves[i - 1],
            Label:    "book",
            CPScore:  moveEvaluation.Score,
            Position: i,
        }
        } else {
            break
        }
    }
    
    analysis.OpeningName = lastOpeningName
    analysis.OpeningEnd = openingEnd
    return openingEnd
}