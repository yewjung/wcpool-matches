package driver

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoMatchesDB() *mongo.Client {
	uri := "mongodb://root:example@matchesb:27017/"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB matchesdb successfully connected and pinged.")
	return client
}

func ConnectMongoPredictionsDB() *mongo.Client {
	uri := "mongodb://root:example@predictionsdb:27017/"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("MongoDB predictionsdb successfully connected and pinged.")
	return client
}

func ConnectMatchesRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "matchesRedis:6379",
		Password: "",
		DB:       0,
	})
}
func ConnectPredictionsRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "predictionsRedis:6379",
		Password: "",
		DB:       0,
	})
}
