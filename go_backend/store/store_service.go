package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// Define the struct wrapper around raw Redis Client
type StorageService struct {
	redisClient *redis.Client
}

// Top level declaration for storeService and Redis context 
var (
	storeService = &StorageService{}
	ctx = context.Background()
)

const CacheDuration = 6 * time.Hour

func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr		: "localhost:6379",
		Password	: "",
		DB			: 0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started succesfully: pong mesg = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

// SaveUrlMapping stores the short -> original mapping. Any error here means
// Redis failed, so it is returned to the caller to handle (e.g. respond 500)
// instead of crashing the whole request with a panic.
func SaveUrlMapping(shortUrl string, originalUrl string) error {
	return storeService.redisClient.Set(shortUrl, originalUrl, CacheDuration).Err()
}

func RetrieveInitialUrl(shortUrl string) (string, error) {
	return storeService.redisClient.Get(shortUrl).Result()
}