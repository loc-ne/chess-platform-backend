package engine

import (
    "encoding/binary"
)

type BitboardGame struct {
    WhitePawns   uint64 `json:"WhitePawns,string"`
    WhiteRooks   uint64 `json:"WhiteRooks,string"`
    WhiteKnights uint64 `json:"WhiteKnights,string"`
    WhiteBishops uint64 `json:"WhiteBishops,string"`
    WhiteQueens  uint64 `json:"WhiteQueens,string"`
    WhiteKing    uint64 `json:"WhiteKing,string"`

    BlackPawns   uint64 `json:"BlackPawns,string"`
    BlackRooks   uint64 `json:"BlackRooks,string"`
    BlackKnights uint64 `json:"BlackKnights,string"`
    BlackBishops uint64 `json:"BlackBishops,string"`
    BlackQueens  uint64 `json:"BlackQueens,string"`
    BlackKing    uint64 `json:"BlackKing,string"`
}

// Constructor function
func NewBitboardGame() *BitboardGame {
    return &BitboardGame{
        WhitePawns:   binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00}),
        WhiteRooks:   binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81}),
        WhiteKnights: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x42}),
        WhiteBishops: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24}),
        WhiteQueens:  binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08}),
        WhiteKing:    binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10}),
        
        BlackPawns:   binary.BigEndian.Uint64([]byte{0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
        BlackRooks:   binary.BigEndian.Uint64([]byte{0x81, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
        BlackKnights: binary.BigEndian.Uint64([]byte{0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
        BlackBishops: binary.BigEndian.Uint64([]byte{0x24, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
        BlackQueens:  binary.BigEndian.Uint64([]byte{0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
        BlackKing:    binary.BigEndian.Uint64([]byte{0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
    }
}