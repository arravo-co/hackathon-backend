package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

var consumer *nsq.Consumer

type MyMetricsChannelMessageHandler struct{}

func (h *MyMetricsChannelMessageHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	fmt.Println(m.Body)

	return nil
}

func init() {
	myMessageHandler := &MyMetricsChannelMessageHandler{}
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("participant_register", "metrics", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(myMessageHandler)
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}
}

func Cleanup() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Gracefully stop the consumer.
	consumer.Stop()
}
