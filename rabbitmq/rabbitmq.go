package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SetupRMQConfig struct {
	Url string
}

func GetRMQChannelWithConn(conn *amqp.Connection) (*amqp.Channel, error) {

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Channel created")
	return ch, err
}

func GetRMQChannelWithURL(opts SetupRMQConfig) (*amqp.Channel, error) {

	conn, err := amqp.Dial(opts.Url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	return ch, err
}

func GetRMQConnWithURL(opts SetupRMQConfig) (*amqp.Connection, error) {

	conn, err := amqp.Dial(opts.Url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to RabbitMQ")
	return conn, err
}
