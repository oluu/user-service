package data

import "github.com/oluu/user-service/contract"

// MockRepository ...
type MockRepository struct {
	GetUserFunc    func(ID string) (*contract.User, error)
	CreateUserFunc func(data *contract.User) error
	UpdateUserFunc func(data *contract.User) error
	DeleteUserFunc func(ID string) error
}

// GetUser ...
func (m *MockRepository) GetUser(ID string) (*contract.User, error) {
	return m.GetUserFunc(ID)
}

// CreateUser ...
func (m *MockRepository) CreateUser(data *contract.User) error {
	return m.CreateUserFunc(data)
}

// UpdateUser ...
func (m *MockRepository) UpdateUser(data *contract.User) error {
	return m.UpdateUserFunc(data)
}

// DeleteUser ...
func (m *MockRepository) DeleteUser(ID string) error {
	return m.DeleteUserFunc(ID)
}
