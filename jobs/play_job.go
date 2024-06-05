package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/queue"
)

type PlayConsumer struct {
}

func StartConsumingPlayQueue() {
	queue, err := queue.GetQueue("play_list")
	if err != nil {
		fmt.Println("Error getting queue")
		return
	}
	err = queue.StartConsuming(1, time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	taskConsumer := &PlayConsumer{}
	str, err := queue.AddConsumer("play_oooo", taskConsumer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)

	for {

	}
}

func (c *PlayConsumer) Consume(d rmq.Delivery) {
	fmt.Println("loading playlist")
	payload := d.Payload()

	payloadStruct := exports.PlayQueuePayload{
		Time: time.Now(),
	}
	err := json.Unmarshal([]byte(payload), &payloadStruct)
	if err != nil {
		fmt.Println(err.Error())
		if err := d.Reject(); err != nil {
			exports.MySugarLogger.Errorln("Failed to reject delivery")
			exports.MySugarLogger.Errorln(err.Error())
		}
		return
	}
	fmt.Println(payloadStruct)
	d.Ack()
	panic("Intentional failure")
}
