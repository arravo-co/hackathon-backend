package db

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	url := config.GetRedisURL()
	opts, err := redis.ParseURL(url)
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	RedisClient = redis.NewClient(opts)
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
	fmt.Println(pong)
}
