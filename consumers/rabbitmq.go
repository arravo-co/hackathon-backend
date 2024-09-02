package consumers

import (
	"github.com/aidarkhanov/nanoid"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQConsumer struct {
	Channel *amqp.Channel
}

func NewRMQConsumerWithChannel(ch *amqp.Channel) *RMQConsumer {
	return &RMQConsumer{
		Channel: ch,
	}
}

func (c *RMQConsumer) DeclareExchange(name string, kind string) error {
	return c.Channel.ExchangeDeclare(name, kind, false, false, false, false, nil)
}

func (c *RMQConsumer) DeclareQueue(name string) (amqp.Queue, error) {
	return c.Channel.QueueDeclare(name, true, false, false, false, nil)
}

// Rare. Should seldom be called

func (c *RMQConsumer) BindQueue(name string, key string, exchange string) error {
	return c.Channel.QueueBind(name, key, exchange, false, nil)
}

func (c *RMQConsumer) ConsumeQueue(q_name string, consumer_name string) (<-chan amqp.Delivery, error) {
	resChan, err := c.Channel.Consume(q_name, consumer_name, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return resChan, err
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
