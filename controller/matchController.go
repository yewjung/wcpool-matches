package controller

import (
	"context"
	"match/authorization"
	"match/models"
	"match/service"
	"match/utils"
	"net/http"
)

type MatchController struct {
	Storage                   models.Storage
	AuthClient                authorization.AuthorizationClient
	matchAndPredictionService *service.MatchAndPredictionService
}

func (mc *MatchController) GetMatchesAndPredictions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, models.MatchRequestDTO{})
		ok, email := mc.checkAuthorization(w, r, dto.Partyid, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			return
		}
		mps := mc.getMatchAndPredictionService()
		result, err := mps.GetMatchesAndPredictions(dto.Matchday, email, dto.Partyid)
		utils.HandleResponse(w, err, result)
	}
}

func (mc *MatchController) AddPrediction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := utils.GetReqBody(r, models.PredictionDTO{})
		ok, email := mc.checkAuthorization(w, r, dto.PartyID, []authorization.Option{authorization.Option_PARTY_ID})
		if !ok {
			return
		}
		mps := mc.getMatchAndPredictionService()
		key := models.MatchEmailPartyKey{
			MatchID: dto.MatchID,
			Email:   email,
			PartyID: dto.PartyID,
		}
		err := mps.AddPrediction(key, models.Prediction{
			GoalA: dto.GoalA,
			GoalB: dto.GoalB,
			Score: dto.Score,
		})
		utils.HandleResponse(w, err, nil)
	}
}

func (mc *MatchController) checkAuthorization(w http.ResponseWriter, r *http.Request, partyid string, options []authorization.Option) (bool, string) {
	verRes, err := mc.AuthClient.VerifyPartyID(context.Background(), &authorization.Verification{
		Token:   r.Header.Get("Authorization"),
		Partyid: partyid,
		Options: options,
	})
	if err != nil {
		utils.HandleResponse(w, err, nil)
		return false, verRes.Email
	}
	if !verRes.Ok {
		utils.SendError(w, http.StatusUnauthorized, models.Error{
			Message: "Unauthorized action",
		})
		return false, verRes.Email
	}
	return true, verRes.Email
}

func (mc *MatchController) getMatchAndPredictionService() *service.MatchAndPredictionService {
	if mc.matchAndPredictionService == nil {
		mc.matchAndPredictionService = &service.MatchAndPredictionService{
			Storage: mc.Storage,
		}
	}
	return mc.matchAndPredictionService
}
