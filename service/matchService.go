package service

import (
	"context"
	"match/models"
	"match/utils"

	"match/repository"
)

type MatchService struct {
	Storage   models.Storage
	matchRepo *repository.MatchMongo
}

func (ms *MatchService) GetMatchesByMatchday(matchday string) []models.Match {
	return utils.GetFromCacheOrFunc(context.Background(), ms.Storage.MatchRedis, models.Key(matchday), ms.getMatchRepo().GetMatchesByMatchday)
}

func (ms *MatchService) getMatchRepo() *repository.MatchMongo {
	if ms.matchRepo == nil {
		ms.matchRepo = &repository.MatchMongo{
			Client: ms.Storage.MatchMongo,
		}
	}
	return ms.matchRepo
}
