package data

import (
	"fmt"
	"log"

	"github.com/oluu/user-service/contract"
)

const dataAccessErr = "Data access request for %s failed"

// Repository ...
type Repository interface {
	GetUser(ID string) (*contract.User, error)
	CreateUser(data *contract.User) error
	UpdateUser(data *contract.User) error
	DeleteUser(ID string) error
}

type repository struct {
	cassandraRepository
}

// InitRepository ...
func InitRepository() Repository {
	sessionManager, err := connectToCassandra()
	if err != nil {
		log.Fatal(err)
	}
	session, err := sessionManager.getSession()
	if err != nil {
		log.Fatal(err)
	}
	return &repository{
		cassandraRepository{session},
	}
}

// GetUser ...
func (r *repository) GetUser(ID string) (*contract.User, error) {
	name := "GetUser"
	iter, err := r.query(name, `SELECT id, email, first_name, last_name FROM user WHERE id=?;`, 0, ID)
	if err != nil {
		log.Println("ERR:", err)
		return nil, fmt.Errorf(dataAccessErr, name)
	}
	if iter.NumRows() == 0 {
		return nil, nil
	}

	record := new(contract.User)
	iter.Scan(&record.ID, &record.Email, &record.FirstName, &record.LastName)
	if err := iter.Close(); err != nil {
		log.Println("ERR:", err)
		return nil, fmt.Errorf(dataAccessErr, name)
	}
	return record, nil
}

// CreateUser ...
func (r *repository) CreateUser(data *contract.User) error {
	name := "CreateUser"
	err := r.execute(name, "INSERT INTO user (id, email, first_name, last_name) VALUES (?, ?, ?, ?);", data.ID, data.Email, data.FirstName, data.LastName)
	if err != nil {
		log.Println("ERR:", err)
		return fmt.Errorf(dataAccessErr, name)
	}
	return nil
}

// UpdateUser ...
func (r *repository) UpdateUser(data *contract.User) error {
	name := "UpdateUser"
	err := r.execute(name, "UPDATE user SET email, first_name=?, last_name=? WHERE id=?;", data.Email, data.FirstName, data.LastName, data.ID)
	if err != nil {
		log.Println("ERR:", err)
		return fmt.Errorf(dataAccessErr, name)
	}
	return nil
}

// DeleteUser ...
func (r *repository) DeleteUser(ID string) error {
	name := "DeleteUser"
	err := r.execute(name, "DELETE FROM user WHERE id=?;", ID)
	if err != nil {
		log.Println("ERR:", err)
		return fmt.Errorf(dataAccessErr, name)
	}
	return nil
}
