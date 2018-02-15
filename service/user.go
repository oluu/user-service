package service

import (
	"encoding/json"
	"log"

	nsq "github.com/nsqio/go-nsq"
	"github.com/shinypotato/user-service/contract"
	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/util"
)

// UserService ...
type UserService struct {
	repository data.Repository
	producer   *nsq.Producer
}

// NewUserService ...
func NewUserService(repository data.Repository, producer *nsq.Producer) *UserService {
	return &UserService{repository: repository, producer: producer}
}

// GetUser ...
func (svc *UserService) GetUser(ID string) (*contract.User, error) {
	user, err := svc.repository.GetUser(ID)
	if err != nil {
		log.Printf("ERR: An error occurred while retrieving user with ID %s. Err: %v\n", ID, err)
		return nil, err
	}
	return user, nil
}

// CreateUser ...
func (svc *UserService) CreateUser(data *contract.User) (string, error) {
	id := util.RandomUUIDString()
	data.ID = id
	err := svc.publishMessage("user-create", data)
	if err != nil {
		log.Printf("ERR: An error occurred while creating user. Err: %v\n", err)
		return "", err
	}
	return id, nil
}

// UpdateUser ...
func (svc *UserService) UpdateUser(data *contract.User) error {
	err := svc.publishMessage("user-update", data)
	if err != nil {
		log.Printf("ERR: An error occurred while updating user with id %s. Err: %v\n", data.ID, err)
		return err
	}
	return nil
}

// DeleteUser ...
func (svc *UserService) DeleteUser(ID string) error {
	err := svc.producer.Publish("user-delete", []byte(ID))
	if err != nil {
		log.Printf("ERR: An error occurred while deleting user with id %s. Err: %v\n", ID, err)
		return err
	}
	return nil
}

func (svc *UserService) publishMessage(topic string, data interface{}) error {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERR: An error occurred while marshaling data. err: %v", err)
		return err
	}
	return svc.producer.Publish(topic, message)
}
