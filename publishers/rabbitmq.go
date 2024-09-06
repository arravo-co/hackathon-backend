package publishers

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var ProducerChannel *amqp.Channel
var ConsumerChannel *amqp.Channel

type RMQPublisher struct {
	Config  *exports.PublisherConfig
	Channel *amqp.Channel
}

func NewRMQPublisherWithChannel(ch *amqp.Channel) *RMQPublisher {
	return &RMQPublisher{
		Channel: ch,
	}
}

func (p *RMQPublisher) DeclareExchange(name string, kind string) error {
	return p.Channel.ExchangeDeclare(name, kind, false, false, false, false, nil)
}

func (p *RMQPublisher) DeclareQueue(name string) (amqp.Queue, error) {
	return p.Channel.QueueDeclare(name, true, false, false, false, nil)
}

// Rare. Should seldom be called
func (p *RMQPublisher) BindQueue(name string, key string, exchange string) error {
	return p.Channel.QueueBind(name, key, exchange, false, nil)
}

func Publish(ProducerChannel *amqp.Channel, exchange string, key string, body []byte) error {
	err := ProducerChannel.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		Body:        body,
		ContentType: "application/json",
		AppId:       "arravo_hackathon",
	})
	return err
}

func (p *RMQPublisher) Publish(ops exports.PublisherConfig, body []byte) error {
	err := Publish(p.Channel, ops.RabbitMQExchange, ops.RabbitMQKey, body)
	//ProducerChannel.NotifyPublish()
	if err != nil {
		return err
	}
	return err
}

func (p *RMQPublisher) DeclareAllExchanges() {
	//c.Channel.
	err := p.Channel.ExchangeDeclare("judge.registered", amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Exchange 'judge.registered' declaration successful")
	//
}

func DeclareAllExchanges(channel *amqp.Channel) {
	//c.Channel.
	err := channel.ExchangeDeclare(exports.InvitationsExchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Printf("Exchange '%s' declaration successful\n", exports.InvitationsExchange)
	}

	err = channel.ExchangeDeclare(exports.JudgesExchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Printf("Exchange '%s' declaration successful\n", exports.JudgesExchange)
	}

	err = channel.ExchangeDeclare(exports.ParticipantsExchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Printf("Exchange '%s' declaration successful\n", exports.ParticipantsExchange)
	}

	err = channel.ExchangeDeclare(exports.UploadJobsExchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Printf("Exchange '%s' declaration successful\n", exports.UploadJobsExchange)
	}
	//
}
