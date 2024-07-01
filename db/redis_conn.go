package db

import (
	"context"
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
	DefaultRedisClient = NewRedisClient()
	return DefaultRedisClient
}

func NewRedisClient() *redis.Client {
	url := config.GetRedisURL()
	opts, err := redis.ParseURL(url)
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	DefaultRedisClient = redis.NewClient(opts)
	pong, err := DefaultRedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	fmt.Println(pong)
	return DefaultRedisClient
}
