package rabbitutils

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

type Publisher struct {
	Config *exports.PublisherConfig
}

func init() {
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

func CreateChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	return ch, err
}

func MustCreateChannel() *amqp.Channel {
	ch, err := CreateChannel()
	if err != nil {
		log.Fatalf("error creating channel: %v", err)
	}
	return ch
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

func Publish(exchange string, key string, body []byte) error {
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

func (p *Publisher) Publish(ops exports.PublisherConfig, body []byte) error {
	err := Publish(ops.RabbitMQExchange, ops.RabbitMQKey, body)
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
	//
}

func ListenToAllQueues() {
	chUploadPicQicJobDelivery, err := ConsumeQueue("upload.profile_picture.cloudinary", GetConsumerTag())
	if err != nil {
		fmt.Println(err.Error())
	}

	chWelcomeEmailToJudgeDelivery, err := ConsumeQueue("send.judge.created.admin.welcome_email", GetConsumerTag())
	if err != nil {
		fmt.Println(err.Error())
	}

	//send.participant.created.welcome_email_verification_email
	chCoordinateParticipantWelcomeVerDelivery, err := ConsumeQueue("send.participant.created.welcome_email_verification_email", GetConsumerTag())
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		select {
		case response := <-chUploadPicQicJobDelivery:
			fmt.Printf(response.ConsumerTag)
			HandlePicUploadConsumption(&response)
		case response := <-chWelcomeEmailToJudgeDelivery:
			HandleSendEmailToJudgeConsumption(&response)
		case response := <-chCoordinateParticipantWelcomeVerDelivery:
			fmt.Printf(response.ConsumerTag)
			//HandleCoordinateParticipantWelcomeVerificationConsumption(&response)
		}
	}
}

func GetConsumerTag() string {

	id := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789", 10))
	return id
}
