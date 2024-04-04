package queue

import (
	"fmt"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/db"
)

var Rmq rmq.Connection

func init() {
	errChan := make(chan error)
	conn, err := rmq.OpenConnectionWithRedisClient("my queue", db.RedisClient, errChan)
	if err != nil {
		fmt.Println(err.Error())
	}
	Rmq = conn
}

func GetQueue(name string) (rmq.Queue, error) {
	return Rmq.OpenQueue(name)
}
