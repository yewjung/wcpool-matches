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

type PredictionMongo struct {
	Client *mongo.Client
}

const databaseName = "Matches"
const collectionName = "Prediction"

func (pm *PredictionMongo) GetPredictionByMatchIDNEmailNPartyID(key models.MatchEmailPartyKey) models.Prediction {
	predictionCollection := pm.Client.Database(databaseName).Collection(collectionName)
	opts := options.FindOne().SetProjection(bson.M{"_id": 0})
	objId, _ := primitive.ObjectIDFromHex(key.String())
	result := predictionCollection.FindOne(context.Background(), bson.M{"_id": objId}, opts)
	pred := models.Prediction{}
	err := result.Decode(&pred) // im not sure if this will work
	if err != nil {
		log.Default().Panic(err)
	}
	pred.MatchID = key.MatchID
	return pred
}

func (pm *PredictionMongo) AddPrediction(key models.MatchEmailPartyKey, prediction models.Prediction) error {
	predictionCollection := pm.Client.Database(databaseName).Collection(collectionName)
	_, err := predictionCollection.InsertOne(context.Background(), bson.M{"_id": key.String(), "prediction": prediction})
	if err != nil {
		log.Default().Panic(err)
	}
	return err
}
