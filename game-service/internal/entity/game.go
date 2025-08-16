package entity

import (
	"time"
)

type PlayerInfo struct {
    UserID   int 	`bson:"userId"`
    Username string `bson:"username"`
	Elo		 int	`bson:"elo"`
}

type Game struct {
    ID         string      `bson:"-"`
	GameID     string      `bson:"gameId"`
    Players    struct {
        White PlayerInfo `bson:"white"`
        Black PlayerInfo `bson:"black"`
    } `bson:"players"`
    Moves       []string  `bson:"moves"`
    Result      string    `bson:"result"`
    CreatedAt   time.Time `bson:"createdAt"`
    TimeControl string    `bson:"timeControl"`
	GameType 	string	  `bson:"gameType"`
    WinnerID    string    `bson:"winnerId"`
    WhiteTimeLeft int       `bson:"whiteTimeLeft"` 
    BlackTimeLeft int       `bson:"blackTimeLeft"`
    Reason           string    `bson:"reason"`
    LastFen         string   `bson:"lastFen"`     
}