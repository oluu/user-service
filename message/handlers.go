package message

import (
	"encoding/json"
	"fmt"
	"log"

	nsq "github.com/nsqio/go-nsq"
	"github.com/shinypotato/user-service/contract"
	"github.com/shinypotato/user-service/data"
)

// InitHandlers ...
func InitHandlers(consumers []*nsq.Consumer, repository data.Repository) {
	consumers[0].AddHandler(CreateUser(repository))
	consumers[1].AddHandler(UpdateUser(repository))
	consumers[2].AddHandler(DeleteUser(repository))
}

// CreateUser ...
func CreateUser(repository data.Repository) nsq.HandlerFunc {
	handler := func(message *nsq.Message) error {
		body := new(contract.User)
		if err := readMessageBody(message, &body); err != nil {
			log.Printf("ERR: error parsing message body, %v\n", err)
			return err
		}
		return repository.CreateUser(body)
	}

	return handler
}

// UpdateUser ...
func UpdateUser(repository data.Repository) nsq.HandlerFunc {
	handler := func(message *nsq.Message) error {
		body := new(contract.User)
		if err := readMessageBody(message, &body); err != nil {
			log.Printf("ERR: error parsing message body, %v\n", err)
			return err
		}
		return repository.UpdateUser(body)
	}
	return handler
}

// DeleteUser ...
func DeleteUser(repository data.Repository) nsq.HandlerFunc {
	handler := func(message *nsq.Message) error {
		id := string(message.Body)
		return repository.DeleteUser(id)
	}
	return handler
}

func readMessageBody(message *nsq.Message, body interface{}) error {
	if err := json.Unmarshal(message.Body, body); err != nil {
		return fmt.Errorf("Failed to read message: %s", err)
	}
	return nil
}
