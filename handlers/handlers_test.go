package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/service"

	"github.com/shinypotato/user-service/contract"
)

type HandleTester func(
	method string,
	params url.Values,
) *httptest.ResponseRecorder

// Given the current test runner and an http.Handler, generate a
// HandleTester which will test its given input against the
// handler.

func generateHandlerTester(t *testing.T, handleFunc http.Handler) HandleTester {
	// Given a method type ("GET", "POST", etc) and
	// parameters, serve the response against the handler and
	// return the ResponseRecorder.

	return func(method string, params url.Values) *httptest.ResponseRecorder {
		req, err := http.NewRequest(method, "", strings.NewReader(params.Encode()))
		if err != nil {
			t.Errorf("%v", err)
		}
		w := httptest.NewRecorder()
		handleFunc.ServeHTTP(w, req)
		return w
	}
}

func getMockService(repository *data.MockRepository) *service.UserService {
	return service.NewUserService(repository)
}

func getMockRepository(getUserErr, createUserErr, updateUserErr, deleteUserErr bool) *data.MockRepository {
	repo := new(data.MockRepository)
	repo.GetUserFunc = func(ID string) (*contract.User, error) {
		if getUserErr {
			return nil, errors.New("getUserErr")
		}
		return &contract.User{
			Email:     "testpotato@shinypotato.com",
			FirstName: "test",
			LastName:  "potato",
		}, nil
	}
	repo.CreateUserFunc = func(data *contract.User) error {
		if createUserErr {
			return errors.New("createUserErr")
		}
		return nil
	}
	repo.UpdateUserFunc = func(data *contract.User) error {
		if updateUserErr {
			return errors.New("updateUserErr")
		}
		return nil
	}
	repo.DeleteUserFunc = func(ID string) error {
		if deleteUserErr {
			return errors.New("deleteUserErr")
		}
		return nil
	}
	return repo
}

func TestGetUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := GetUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodGet, url.Values{})
	if expected, actual := http.StatusOK, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
	expected, _ := json.Marshal(&contract.User{
		Email:     "testpotato@shinypotato.com",
		FirstName: "test",
		LastName:  "potato",
	})
	if actual := w.Body.String(); string(expected) != actual {
		t.Errorf("Wrong body returned: expected %v, actual %v", expected, actual)
	}
}

func TestCreateUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := CreateUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodPost, url.Values{})
	if expected, actual := http.StatusCreated, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}

func TestUpdateUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := UpdateUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodPut, url.Values{})

	if expected, actual := http.StatusAccepted, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}

func TestDeleteUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := DeleteUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodDelete, url.Values{})
	if expected, actual := http.StatusAccepted, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}
