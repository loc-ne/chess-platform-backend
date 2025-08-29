package usecase

import ("math")

const (
    BEST_MOVE_THRESHOLD     = 0.3  
    EXCELLENT_MOVE_THRESHOLD = 0.6   
    GOOD_MOVE_THRESHOLD     = 0.8   
    INACCURACY_THRESHOLD    = 1   
    MISTAKE_THRESHOLD       = 1.5  
)

func classifyMove(prevEval, currentEval MoveEvaluation, isFirstMate *bool, mateInMoves *int, moveIndex, totalMoves int, inputMoves []string) string {
    // Update mate tracking variables exactly as in original
    if currentEval.IsMate {
        if !*isFirstMate {
            *isFirstMate = true
        }
        *mateInMoves = currentEval.MateIn
        
        // Original mate classification logic
        if *isFirstMate {
            isWhiteInMate := determineWhiteInMate(currentEval)
            
            if isWhiteInMate != currentEval.IsWhiteMove && moveIndex+1 != totalMoves {
                if prevEval.Score > 0 {
                    return "blunder"
                } else if prevEval.Score < 0 {
                    if prevEval.Score <= -1.0 {
                        return "mistake"
                    } else {
                        return "blunder"
                    }
                }
            } else if moveIndex+1 == totalMoves && currentEval.MateIn == 0 {
                return "best"
            } else if *mateInMoves > currentEval.MateIn {
                return "best"
            } else {
                return "excellent"
            }
        }
    } else {
        *isFirstMate = false
    }
    
     return classifyNormalMoveWithBrilliant(prevEval.Score, currentEval.Score, inputMoves, moveIndex)
}

func classifyNormalMoveWithBrilliant(prevScore, currentScore float64, inputMoves []string, moveIndex int) string {
    moveLabel := classifyNormalMove(prevScore, currentScore)
    
    // Check for brilliant move only if it's classified as "best move"
    if moveLabel == "best" || moveLabel == "excellent" || moveLabel == "good" {
        if isBrilliantMove(inputMoves, moveIndex) {
            return "brilliant"
        }
    }
    
    return moveLabel
}

func determineWhiteInMate(eval MoveEvaluation) bool {
    if eval.MateIn < 0 {
        return !eval.IsWhiteMove
    } else {
        return eval.IsWhiteMove
    }
}

func classifyNormalMove(prevScore, currentScore float64) string {
    scoreDiff := math.Abs(currentScore - prevScore)
    
    // Base classification
    var moveLabel string
    if scoreDiff <= BEST_MOVE_THRESHOLD {
        moveLabel = "best"
    } else if scoreDiff <= EXCELLENT_MOVE_THRESHOLD {
        moveLabel = "excellent"
    } else if scoreDiff <= GOOD_MOVE_THRESHOLD {
        moveLabel = "good"
    } else if scoreDiff <= INACCURACY_THRESHOLD {
        moveLabel = "inaccuracy"
    } else if scoreDiff > INACCURACY_THRESHOLD {
        moveLabel = "mistake"
    } else if scoreDiff > MISTAKE_THRESHOLD {
        moveLabel = "blunder"
    }

    // Original advantage flip detection
    if (prevScore > 0.0 && currentScore < 0.0) || (prevScore < 0.0 && currentScore > 0.0) {
        if scoreDiff <= 1.5 {
            moveLabel = "mistake"
        } else {
            moveLabel = "blunder"
        }
    }
    
    return moveLabel
}