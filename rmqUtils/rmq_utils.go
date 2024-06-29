package rmqUtils

import (
	"fmt"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/db"
)

var DefaultRmq rmq.Connection
var ErrCh chan error

func SetupDefaultQueue() {
	var err error
	DefaultRmq, err = GetDefaultQueue()
	if err != nil {
		panic(err.Error())
	}
}
func GetDefaultQueue() (rmq.Connection, error) {
	ErrCh = make(chan error)
	conn, err := rmq.OpenConnectionWithRedisClient("my queue", db.DefaultRedisClient, ErrCh)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	DefaultRmq = conn
	return DefaultRmq, nil
}

func GetQueue(name string) (rmq.Queue, error) {
	return DefaultRmq.OpenQueue(name)
}
