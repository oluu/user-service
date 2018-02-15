package message

import (
	"log"

	nsq "github.com/nsqio/go-nsq"
)

// InitMessaging ...
func InitMessaging(topics []string) (*nsq.Producer, *nsq.Consumer) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer, err := nsq.NewConsumer("user", "default", config)
	if err != nil {
		log.Fatal(err)
	}
	consumer.ChangeMaxInFlight(200)
	return producer, consumer
}
