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