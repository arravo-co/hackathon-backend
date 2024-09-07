package db

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/redis/go-redis/v9"
)

var DefaultRedisClient *redis.Client

func SetupRedis() {
	if DefaultRedisClient != nil {
		return
	}
	DefaultRedisClient = GetDefaultRedisClient()
}

func GetDefaultRedisClient() *redis.Client {
	if DefaultRedisClient != nil {
		return DefaultRedisClient
	}
	DefaultRedisClient = NewRedisDefaultClient()
	return DefaultRedisClient
}

func NewRedisDefaultClient() *redis.Client {
	url := config.GetRedisURL()
	opts, err := redis.ParseURL(url)
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	DefaultRedisClient = redis.NewClient(opts)
	return DefaultRedisClient
}
