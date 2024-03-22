package db

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	url, _ := config.GetRedisURL()
	opts, err := redis.ParseURL(url)
	if err != nil {
		fmt.Printf("%v", err)
	}
	RedisClient = redis.NewClient(opts)
}
