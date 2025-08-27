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

