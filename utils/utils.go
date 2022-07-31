package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"match/models"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
)

func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func SendServerErrorIfErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		errorMessage := models.Error{Message: "Server Error"}
		SendError(w, http.StatusInternalServerError, errorMessage)
		return true
	}
	return false
}

func HandleResponse(w http.ResponseWriter, err error, data interface{}) {
	if err != nil {
		log.Default().Panic(err)
		errorMessage := models.Error{Message: "Server Error"}
		SendError(w, http.StatusInternalServerError, errorMessage)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	SendSuccess(w, data)
}

func RemoveFromArray(arr []int, item int) []int {
	for i, v := range arr {
		if v == item {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}
	return arr
}

func GetIntVar(r *http.Request, key string) int {
	value := mux.Vars(r)[key]
	intVal, _ := strconv.Atoi(value)
	return intVal
}

func GetReqBody[T any](r *http.Request, data T) T {
	json.NewDecoder(r.Body).Decode(&data)
	return data
}

func Map[T any, R any](slice []T, f func(T) R) []R {
	newSlice := make([]R, len(slice))
	for i, t := range slice {
		newSlice[i] = f(t)
	}
	return newSlice
}

func GetFromCacheOrFunc[K fmt.Stringer, T any](ctx context.Context, cache *redis.Client, key K, f func(K) T) T {
	result, err := cache.Get(context.Background(), key.String()).Result()
	if err == nil {
		value := new(T)
		json.Unmarshal([]byte(result), &value)
		return *value
	}
	freshValue := f(key)
	encodedValue, err := json.Marshal(freshValue)
	if err != nil {
		log.Default().Panic(err)
		return freshValue
	}
	cache.Set(ctx, key.String(), encodedValue, 0)
	return freshValue
}

func WriteAroundCache[K fmt.Stringer, T any](ctx context.Context, cache *redis.Client, key K, value T, f func(K, T) error) error {
	err := cache.Del(ctx, key.String()).Err()
	if err != nil {
		log.Default().Panic(err)
	}
	return f(key, value)

}
