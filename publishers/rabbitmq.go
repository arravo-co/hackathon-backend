package publishers

import (
	"context"

	"github.com/aidarkhanov/nanoid"
	"github.com/arravoco/hackathon_backend/exports"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var ProducerChannel *amqp.Channel
var ConsumerChannel *amqp.Channel

type RMQPublisher struct {
	//Config  *exports.PublisherConfig
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
	//ProducerChannel.NotifyPublish()
	if err != nil {
		return err
	}
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

/*
func DeclareAllQueues() {
	q, err := DeclareQueue("upload.profile_picture.cloudinary")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(q)

	sendJudgeCreatedAdminWelcomeEmailQueue, err := DeclareQueue("send.judge.created.admin.welcome_email")
	if err != nil {
		log.Fatalln("Error: ", err.Error())
	}
	log.Println(sendJudgeCreatedAdminWelcomeEmailQueue)

	sendParticipantCreatedWelcomeEmailQueue, err := DeclareQueue("send.participant.created.welcome_email_verification_email")
	if err != nil {
		log.Fatalln("Error: ", err.Error())
	}
	log.Println(sendParticipantCreatedWelcomeEmailQueue)

	uploadSolutionPicQueue, err := DeclareQueue("upload.solution_picture.cloudinary")
	if err != nil {
		log.Fatalln("Error: ", err.Error())
	}
	log.Println(uploadSolutionPicQueue)
	//
}*/

func GetConsumerTag() string {
	id := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789", 10))
	return id
}
