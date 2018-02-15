package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/shinypotato/user-service/contract"
	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/service"
	"github.com/shinypotato/user-service/util"
)

type HandleTester func(
	method string,
	path string,
	params map[string]string,
) *httptest.ResponseRecorder

const (
	envProtocol = "PROTOCOL"
	envHost     = "HOST"
	envPort     = "PORT"
)

func generateHandlerTester(t *testing.T, handleFunc http.Handler) HandleTester {
	return func(method, path string, values map[string]string) *httptest.ResponseRecorder {
		jsonValue, _ := json.Marshal(values)
		req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonValue))
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

func buildURL(path string, queryPairs ...string) string {
	baseURL := fmt.Sprintf("%s://%s:%d", util.GetEnvStringOrDefault(envProtocol, "http"), util.GetEnvString(envHost), util.GetEnvInt(envPort))
	if queryPairs != nil && len(queryPairs) > 0 {
		qs := url.Values{}
		for i := 0; i+1 < len(queryPairs); i += 2 {
			qs.Set(queryPairs[i], queryPairs[i+1])
		}
		return baseURL + path + "?" + qs.Encode()
	}
	return baseURL + path
}

func TestGetUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := routeTraffic(svc) //GetUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodGet, buildURL(contract.GetUser, "id", "123"), nil)
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
	handler := routeTraffic(svc) //CreateUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodPost, buildURL(contract.PostUser), nil)
	if expected, actual := http.StatusCreated, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}

func TestUpdateUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := routeTraffic(svc) //UpdateUser(svc)
	test := generateHandlerTester(t, handler)
	values := map[string]string{"email": "testpotato@shinypotato.com"}
	w := test(http.MethodPut, buildURL(contract.PutUser, "id", "123"), values)
	if expected, actual := http.StatusAccepted, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}

func TestDeleteUser(t *testing.T) {
	svc := getMockService(getMockRepository(false, false, false, false))
	handler := routeTraffic(svc) //DeleteUser(svc)
	test := generateHandlerTester(t, handler)
	w := test(http.MethodDelete, buildURL(contract.DeleteUser, "id", "123"), nil)
	if expected, actual := http.StatusAccepted, w.Code; expected != actual {
		t.Errorf("Wrong status code returned: expected %v, actual %v", expected, actual)
	}
}
