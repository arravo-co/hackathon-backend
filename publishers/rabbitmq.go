package publishers

import (
	"context"
	"fmt"
	"log"

	"github.com/aidarkhanov/nanoid"
	"github.com/arravoco/hackathon_backend/config"
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

type RMQConsumer struct {
	Channel *amqp.Channel
}

type SetupRMQConfig struct {
	Url string
}

func NewPublisherWithChannel(ch *amqp.Channel) *RMQPublisher {
	return &RMQPublisher{
		Channel: ch,
	}
}

func SetupDefaultRMQ() {
	var err error
	rabbitMQURL := config.GetRabbitMQURL()
	fmt.Println("Rabbit URL: ", rabbitMQURL)
	Conn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to RabbitMQ")
	ConsumerChannel, err = CreateChannel()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ConsumerChannel to RabbitMQ")

	ProducerChannel, err = CreateChannel()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ProducerChannel to RabbitMQ")
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

func CreateChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	return ch, err
}

func (p *RMQPublisher) DeclareExchange(name string, kind string) error {
	return p.Channel.ExchangeDeclare(name, kind, false, false, false, false, nil)
}

func (p *RMQPublisher) DeclareQueue(name string) (amqp.Queue, error) {
	return p.Channel.QueueDeclare(name, true, false, false, false, nil)
}

func (c *RMQConsumer) DeclareQueue(name string) (amqp.Queue, error) {
	return c.Channel.QueueDeclare(name, true, false, false, false, nil)
}

func DeclareExchange(name string, kind string) error {
	ch := ConsumerChannel
	return ch.ExchangeDeclare(name, kind, false, false, false, false, nil)
}

func DeclareQueue(name string) (amqp.Queue, error) {
	ch := ConsumerChannel
	return ch.QueueDeclare(name, true, false, false, false, nil)
}

func BindQueue(name string, key string, exchange string) error {
	ch := ConsumerChannel
	return ch.QueueBind(name, key, exchange, false, nil)
}

func ConsumeQueue(q_name string, consumer_name string) (<-chan amqp.Delivery, error) {
	resChan, err := ConsumerChannel.Consume(q_name, consumer_name, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return resChan, err
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

func DeclareAllPublisherExchanges() {
	//DeclareExchange("upload_pic_cloudinary", "")
}

func DeclareAllConsumerExchanges() {
	//DeclareExchange("upload_pic_cloudinary", "")
}

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
}

func GetConsumerTag() string {

	id := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789", 10))
	return id
}
