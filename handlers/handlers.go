package handlers

import (
	"net/http"

	"github.com/shinypotato/user-service/contract"
	"github.com/shinypotato/user-service/service"
)

// RegisterHandlers ...
func RegisterHandlers(svc *service.UserService) {
	// http.HandleFunc("/user/{id}", routeTraffic(svc))
	http.HandleFunc("/user", routeTraffic(svc))
}

func routeTraffic(svc *service.UserService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetUser(svc, w, r)
		}
		if r.Method == http.MethodPost {
			CreateUser(svc, w, r)
		}
		if r.Method == http.MethodPut {
			UpdateUser(svc, w, r)
		}
		if r.Method == http.MethodDelete {
			DeleteUser(svc, w, r)
		}
	})
}

// GetUser ...
func GetUser(svc *service.UserService, w http.ResponseWriter, r *http.Request) {
	ID, err := getRequiredVar(w, r, "id")
	if err != nil {
		return
	}
	result, err := svc.GetUser(ID)
	queryCompleted(w, result, err)
}

// CreateUser ...
func CreateUser(svc *service.UserService, w http.ResponseWriter, r *http.Request) {
	req := new(contract.User)
	err := readRequest(w, r, req)
	if err != nil {
		return
	}
	id, err := svc.CreateUser(req)
	w.Header().Set("Location", id)
	accepted(w, err)
}

// UpdateUser ...
func UpdateUser(svc *service.UserService, w http.ResponseWriter, r *http.Request) {
	ID, err := getRequiredVar(w, r, "id")
	if err != nil {
		return
	}
	req := new(contract.User)
	err = readRequest(w, r, req)
	if err != nil {
		return
	}
	req.ID = ID
	err = svc.UpdateUser(req)
	accepted(w, err)
}

// DeleteUser ...
func DeleteUser(svc *service.UserService, w http.ResponseWriter, r *http.Request) {
	ID, err := getRequiredVar(w, r, "id")
	if err != nil {
		return
	}
	err = svc.DeleteUser(ID)
	accepted(w, err)
}
