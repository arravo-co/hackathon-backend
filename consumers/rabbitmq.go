package consumers

import (
	"fmt"
	"time"

	"github.com/aidarkhanov/nanoid"
	consumerhandlers "github.com/arravoco/hackathon_backend/consumer_handlers"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateRMQConsumerOpts struct {
	Channel       *amqp.Channel
	Logger        *zap.Logger
	RMQConnection *amqp.Connection
}
type RMQConsumer struct {
	Channel       *amqp.Channel
	Logger        *zap.Logger
	RMQConnection *amqp.Connection
}

func NewRMQConsumerWithChannel(opts CreateRMQConsumerOpts) *RMQConsumer {
	return &RMQConsumer{
		Channel:       opts.Channel,
		Logger:        opts.Logger,
		RMQConnection: opts.RMQConnection,
	}
}

func (c *RMQConsumer) SetupChannel() error {
	if c.Channel != nil && !c.Channel.IsClosed() {

	}
	ch, err := c.RMQConnection.Channel()
	if err != nil {
		return err
	}
	c.Channel = ch
	return nil
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

/**/
func (c *RMQConsumer) DeclareAllQueues() {
	q, err := c.DeclareQueue("judge.send.welcome_email") //upload.profile_picture.cloudinary
	if err != nil {
		c.Logger.Fatal(err.Error())
	}
	c.Logger.Sugar().Info(q)
	fmt.Printf("Queue name: %s\n", q.Name)
	fmt.Println("BindToJudgeCreateExchanged called")
	for {
		err := c.BindQueue(q.Name,
			"judge.send.welcome_email", "judge.registered")
		if err != nil {
			fmt.Println(err.Error())
			if c.Channel.IsClosed() {
				err = c.SetupChannel()
				if err != nil {
					break
				}
				continue
			}
			time.Sleep(time.Second * 10)
			continue
		} else {
			fmt.Println("JudgeCreated bound")
		}

		chWelcomeEmailToJudgeDelivery, err := c.ConsumeQueue("judge.send.welcome_email", GetConsumerTag())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(chWelcomeEmailToJudgeDelivery)
		}
		for v := range chWelcomeEmailToJudgeDelivery {
			fmt.Println("responsekkkkkkkkkkkkkkkkkkkkkkkkkkkffffffffffffffkkkkkkkkkkkkk")
			consumerhandlers.HandleSendWelcomeAndEmailVerificationEmailToJudgeConsumption(v.Body)
			v.Ack(false)
		}
	}
}

func (c *RMQConsumer) DeclareAllQueuesParameterized(handler func([]byte) error, exchange, queue_name, key string) error {
	q, err := c.DeclareQueue(queue_name) //upload.profile_picture.cloudinary
	if err != nil {
		c.Logger.Fatal(err.Error())
		return err
	}

	c.Logger.Sugar().Infof("\nAttempting to bind queue %s to exchange %s with key %s\n", queue_name, exchange, key)
	for {
		err := c.BindQueue(q.Name, key, exchange)
		if err != nil {
			fmt.Printf("Failed to bind: %s\n", err.Error())
			if c.Channel.IsClosed() {
				time.Sleep(time.Second * 10)
				err = c.SetupChannel()
				if err != nil {
					break
				}
				continue
			}
			time.Sleep(time.Second * 10)
			continue
		} else {
			fmt.Println("Exchange created bound")
		}

		delivery, err := c.ConsumeQueue(queue_name, GetConsumerTag())
		if err != nil {
			fmt.Println("\nFailed to consume: ", err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		for v := range delivery {
			err := handler(v.Body)
			if err != nil {
				v.Nack(false, false)
				continue
			}
			v.Ack(false)
		}
	}
	return nil
}

func GetConsumerTag() string {
	id := nanoid.Must(nanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456456789", 10))
	return id
}
