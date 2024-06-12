package publish

import (
	"log"

	"github.com/arravoco/hackathon_backend/exports"
)

var rabPub *Publisher

type PublisherInterface interface {
	Publish(ops exports.PublisherConfig, body []byte) error
}

type Publisher struct {
	Publisher PublisherInterface
}

func GetPublisher() *Publisher {
	return rabPub
}

func SetPublisher(pub PublisherInterface) {
	rabPub = &Publisher{
		Publisher: pub,
	}
}

func Publish(ops *exports.PublisherConfig, body []byte) error {
	pub := GetPublisher()
	if pub == nil {
		log.Fatalln("Error: ", "Publisher is nil")
	}
	return pub.Publisher.Publish(*ops, body)
}
