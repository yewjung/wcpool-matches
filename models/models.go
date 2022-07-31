package models

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	MatchMongo      *mongo.Client
	PredictionMongo *mongo.Client
	MatchRedis      *redis.Client
	PredictionRedis *redis.Client
}

type Error struct {
	Message string `json:"message"`
}

type MatchRequestDTO struct {
	Matchday string `json:"matchday"`
	Partyid  string `json:"partyid"`
}

type MatchAndPrediction struct {
	MatchID   string
	TeamA     string
	TeamB     string
	GoalA     int
	GoalB     int
	GoalAPred int
	GoalBPred int
	Date      time.Time
	Score     int
}

type Match struct {
	MatchID string
	TeamA   string
	TeamB   string
	GoalA   int
	GoalB   int
	Date    time.Time
}

type Prediction struct {
	MatchID string `json:"matchId"`
	GoalA   int    `json:"goalA" bson:"GoalA"`
	GoalB   int    `json:"goalB" bson:"GoalB"`
	Score   int    `json:"score" bson:"Score"`
}

type PredictionDTO struct {
	Prediction
	PartyID string `json:"partyId"`
}

type Key string

func (k Key) String() string {
	return string(k)
}

type MatchEmailPartyKey struct {
	MatchID string
	Email   string
	PartyID string
}

func (key MatchEmailPartyKey) String() string {
	return strings.Join([]string{key.MatchID, key.Email, key.PartyID}, "$")
}
