package service

import (
	"context"
	"match/models"
	"match/utils"
)

type MatchAndPredictionService struct {
	Storage           models.Storage
	matchService      *MatchService
	predictionService *PredictionService
}

func (mps *MatchAndPredictionService) GetMatchesAndPredictions(matchday string, email string, partyid string) (map[string][]models.MatchAndPrediction, error) {
	// Get Matches by matchday
	matches := mps.getMatchService().GetMatchesByMatchday(matchday)

	// For each matches, get Predictions by email$partyid

	matchIDMap := make(map[string]models.Match, len(matches))
	matchids := make([]string, len(matches))
	for i, match := range matches {
		matchIDMap[match.MatchID] = match
		matchids[i] = match.MatchID
	}

	preds, err := mps.getPredictionService().GetPredictionsByMatchIDsAndEmailAndPartyID(
		matchids,
		email,
		partyid,
	)

	if err != nil {
		return nil, err
	}

	results := []models.MatchAndPrediction{}
	for _, pred := range preds {
		match := matchIDMap[pred.MatchID]
		matchNPred := models.MatchAndPrediction{
			MatchID:   match.MatchID,
			TeamA:     match.TeamA,
			TeamB:     match.TeamB,
			GoalA:     match.GoalA,
			GoalB:     match.GoalB,
			GoalAPred: pred.GoalA,
			GoalBPred: pred.GoalB,
			Date:      match.Date,
			Score:     pred.Score,
		}
		results = append(results, matchNPred)
	}

	return map[string][]models.MatchAndPrediction{
		matchday: results,
	}, nil
}

func (mps *MatchAndPredictionService) AddPrediction(key models.MatchEmailPartyKey, prediction models.Prediction) error {
	ps := mps.getPredictionService()
	return utils.WriteAroundCache(context.Background(), mps.Storage.PredictionRedis, key, prediction, ps.predMongo.AddPrediction)
}

func (mps *MatchAndPredictionService) getMatchService() *MatchService {
	if mps.matchService == nil {
		mps.matchService = &MatchService{
			Storage: mps.Storage,
		}
	}
	return mps.matchService
}
func (mps *MatchAndPredictionService) getPredictionService() *PredictionService {
	if mps.matchService == nil {
		mps.predictionService = &PredictionService{
			MongoDB: mps.Storage.MatchMongo,
			Cache:   mps.Storage.MatchRedis,
		}
	}
	return mps.predictionService
}
