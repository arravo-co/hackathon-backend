package producer

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

var Producer *nsq.Producer

func init() {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	Producer = producer
}

func Publish(topic string, message []byte) error {
	err := Producer.Publish(topic, message)
	if err != nil {

		fmt.Println(err.Error())
		return err
	}
	return nil
}

func PublishAsync(topic string, message []byte, ch chan *nsq.ProducerTransaction) error {
	err := Producer.PublishAsync(topic, message, ch)
	return err
}
