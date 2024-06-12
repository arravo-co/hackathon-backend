package consume

type ConsumerInterface interface {
	Consume(ch chan<- interface{})
}

type Consumer struct {
	Queue interface{}
}

func SetConsumer() {

}

func (c *Consumer) Consume() {

}
