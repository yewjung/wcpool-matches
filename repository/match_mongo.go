package repository

import (
	"context"
	"log"
	"match/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchMongo struct {
	Client *mongo.Client
}

func (mm *MatchMongo) GetMatchesByMatchday(matchday models.Key) []models.Match {
	matchesCollection := mm.Client.Database("Matches").Collection("Matches")
	opts := options.FindOne().SetProjection(bson.M{"_id": 0})
	objId, _ := primitive.ObjectIDFromHex(matchday.String())
	result := matchesCollection.FindOne(context.Background(), bson.M{"_id": objId}, opts)
	matches := []models.Match{}
	err := result.Decode(&matches)
	if err != nil {
		log.Default().Panic(err)
	}
	return matches
}
