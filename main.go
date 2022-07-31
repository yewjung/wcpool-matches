package main

import (
	"log"
	"match/authorization"
	"match/controller"
	"match/driver"
	"match/models"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	storage := models.Storage{
		MatchMongo:      driver.ConnectMongoMatchesDB(),
		PredictionMongo: driver.ConnectMongoPredictionsDB(),
		MatchRedis:      driver.ConnectMatchesRedis(),
		PredictionRedis: driver.ConnectPredictionsRedis(),
	}
	matchController := controller.MatchController{
		Storage:    storage,
		AuthClient: getSecurityGrpcClient(),
	}
	router := mux.NewRouter()
	router.HandleFunc("/matches", matchController.GetMatchesAndPredictions()).Methods("POST")
	router.HandleFunc("/prediction", matchController.AddPrediction()).Methods("POST")

}

func getSecurityGrpcClient() authorization.AuthorizationClient {
	conn, err := grpc.Dial("security:8085")
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return authorization.NewAuthorizationClient(conn)
}
