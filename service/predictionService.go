package service

import (
	"context"
	"fmt"
	"match/models"
	"match/repository"
	"match/utils"
	"sync"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type PredictionService struct {
	MongoDB   *mongo.Client
	Cache     *redis.Client
	predMongo *repository.PredictionMongo
}

func (ps *PredictionService) GetPredictionsByMatchIDsAndEmailAndPartyID(matchids []string, email string, partyid string) ([]models.Prediction, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(matchids))
	predictionChannel := make(chan models.Prediction, len(matchids))
	predictions := make(chan []models.Prediction)
	for _, matchid := range matchids {
		key := models.MatchEmailPartyKey{
			MatchID: matchid,
			Email:   email,
			PartyID: partyid,
		}
		go ps.producePredictions(predictionChannel, &wg, key)
	}
	go ps.consumePredictions(predictionChannel, predictions)
	wg.Wait()
	close(predictionChannel)

	return <-predictions, nil
}

func (ps *PredictionService) getPredMongo() *repository.PredictionMongo {
	if ps.predMongo != nil {
		return ps.predMongo
	}
	ps.predMongo = &repository.PredictionMongo{
		Client: ps.MongoDB,
	}
	return ps.predMongo
}

func (ps *PredictionService) producePredictions(predChannel chan<- models.Prediction, wg *sync.WaitGroup, key models.MatchEmailPartyKey) {
	predChannel <- utils.GetFromCacheOrFunc(context.Background(), ps.Cache, key, ps.getPredMongo().GetPredictionByMatchIDNEmailNPartyID)
	wg.Done()
}

func (ps *PredictionService) consumePredictions(predChannel <-chan models.Prediction, predictions chan<- []models.Prediction) {
	preds := []models.Prediction{}
	for pred := range predChannel {
		preds = append(preds, pred)
	}
	predictions <- preds
}

func constructUserPartyKey(email, partyid string) string {
	return fmt.Sprintf("%s$%s", email, partyid)
}
