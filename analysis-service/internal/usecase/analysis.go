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

func analyzeGamePhase(inputMoves []string, openingEnd int, engine *StockfishEngine, analysis *GameAnalysis) error {
    prevEval, err := getInitialEvaluation(inputMoves, openingEnd, engine)
    if err != nil {
        return err
    }
    
    isFirstMate := false
    var mateInMoves int

    for i := openingEnd; i < len(inputMoves); i++ {
        currentEval, err := getCurrentEvaluation(inputMoves, i, engine)
        if err != nil {
            continue
        }
        
        moveLabel := classifyMove(prevEval, currentEval, &isFirstMate, &mateInMoves, i, len(inputMoves), inputMoves)
        
        analysis.Moves[i] = MoveAnalysis{
            Move:     inputMoves[i],
            Label:    moveLabel,
            CPScore:  currentEval.Score,
            Position: i + 1,
        }
        
        fmt.Printf("Move %d: %s - %s (Score: %.2f → %.2f, Diff: %.2f)\n", 
            i+1, inputMoves[i], moveLabel, prevEval.Score, currentEval.Score, 
            math.Abs(currentEval.Score - prevEval.Score))
        
        prevEval = currentEval
    }

    return nil
}

func getInitialEvaluation(inputMoves []string, openingEnd int, engine *StockfishEngine) (MoveEvaluation, error) {
    openingMoves := inputMoves[:openingEnd]
    engine.SetPosition(openingMoves)
    
    score, err, isMate, mateIn := engine.GetCurrentScore()
    if err != nil {
        score = 0
    }
    
    isWhiteMove := openingEnd%2 == 0
    return createEvaluation(score, isMate, mateIn, isWhiteMove), nil
}

func getCurrentEvaluation(inputMoves []string, moveIndex int, engine *StockfishEngine) (MoveEvaluation, error) {
    nextMoves := inputMoves[:moveIndex+1]
    engine.SetPosition(nextMoves)
    
    score, err, isMate, mateIn := engine.GetCurrentScore()
    if err != nil {
        score = 0
    }
    
    isWhiteMove := (moveIndex+1)%2 == 0
    return createEvaluation(score, isMate, mateIn, isWhiteMove), nil
}

func createEvaluation(score float64, isMate bool, mateIn int, isWhiteMove bool) MoveEvaluation {
    eval := MoveEvaluation{
        Score:       score,
        IsMate:      isMate,
        MateIn:      mateIn,
        IsWhiteMove: isWhiteMove,
    }
    
    if !isWhiteMove {
        eval.Score = -score
        if isMate {
            eval.MateIn = -mateIn
        }
    }
    
    return eval
}