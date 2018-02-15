package message

import (
	"log"

	nsq "github.com/nsqio/go-nsq"
)

// InitMessaging ...
func InitMessaging(topics []string) (*nsq.Producer, []*nsq.Consumer) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	consumers := make([]*nsq.Consumer, len(topics))
	for i, topic := range topics {
		consumers[i], err = nsq.NewConsumer(topic, "default", config)
		if err != nil {
			log.Fatal(err)
		}
		consumers[i].ChangeMaxInFlight(100)
	}
	return producer, consumers
}
