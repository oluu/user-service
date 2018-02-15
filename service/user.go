package service

import (
	"log"

	"github.com/shinypotato/user-service/contract"
	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/util"
)

// UserService ...
type UserService struct {
	repository data.Repository
}

// NewUserService ...
func NewUserService(repository data.Repository) *UserService {
	return &UserService{repository: repository}
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
	err := svc.repository.CreateUser(data)
	if err != nil {
		log.Printf("ERR: An error occurred while creating user. Err: %v\n", err)
		return "", err
	}
	return id, nil
}

// UpdateUser ...
func (svc *UserService) UpdateUser(data *contract.User) error {
	err := svc.repository.CreateUser(data)
	if err != nil {
		log.Printf("ERR: An error occurred while updating user with id %s. Err: %v\n", data.ID, err)
		return err
	}
	return nil
}

// DeleteUser ...
func (svc *UserService) DeleteUser(ID string) error {
	err := svc.repository.DeleteUser(ID)
	if err != nil {
		log.Printf("ERR: An error occurred while deleting user with id %s. Err: %v\n", ID, err)
		return err
	}
	return nil
}
