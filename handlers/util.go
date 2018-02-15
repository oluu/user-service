package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

const (
	contentType     = "Content-Type"
	contentTypeJSON = "application/json"
)

func writeJSON(w http.ResponseWriter, o interface{}, status int) {
	b, err := json.Marshal(o)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Add(contentType, contentTypeJSON)
	w.Write(b)
}

func getVar(r *http.Request, varName string) string {
	v := mux.Vars(r)[varName]
	if v == "" {
		q := r.URL.Query()
		return q.Get(varName)
	}
	return v
}

func getRequiredVar(w http.ResponseWriter, r *http.Request, varName string) (string, error) {
	m := mux.Vars(r)[varName]
	if m == "" {
		q := r.URL.Query()
		v := q.Get(varName)
		if v == "" || v == "undefined" {
			err := fmt.Sprintf("Required param not provided: %s\n", varName)
			log.Println("ERR:", err)
			badRequest(w, err)
			return "", fmt.Errorf(err)
		}
		return v, nil
	}
	return m, nil
}

func readRequest(w http.ResponseWriter, r *http.Request, body interface{}) error {
	ct := r.Header.Get(contentType)
	if ct != "" && ct != contentTypeJSON {
		unsupportedMedia(w)
		return errors.New(http.StatusText(http.StatusUnsupportedMediaType))
	}
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		errStr := fmt.Sprintf("Error while decoding request body: %v\n", err)
		log.Println("ERR:", errStr)
		unprocessibleEntity(w, errStr)
		return err
	}
	return nil
}

func accepted(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("ERR: Internal Server Error: %v\n", err)
		internalServerError(w, err)
	}
	w.WriteHeader(http.StatusAccepted)
}

func created(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("ERR: Internal Server Error: %v\n", err)
		internalServerError(w, err)
	}
	w.WriteHeader(http.StatusCreated)
}

func acceptedWithBody(w http.ResponseWriter, body interface{}, err error) {
	if err != nil {
		internalServerError(w, err)
		return
	}
	iv := reflect.ValueOf(body)
	if body == nil || iv.Kind().String() == "ptr" && iv.IsNil() {
		log.Println("ERR: Resource not found")
		notFound(w)
		return
	}
	writeJSON(w, body, http.StatusAccepted)
}

func queryCompleted(w http.ResponseWriter, body interface{}, err error) {
	if err != nil {
		internalServerError(w, err)
		return
	}
	iv := reflect.ValueOf(body)
	if body == nil || iv.Kind().String() == "ptr" && iv.IsNil() {
		notFound(w)
		return
	}
	writeJSON(w, body, http.StatusOK)
}

func badRequest(w http.ResponseWriter, err string) {
	log.Printf("ERR: Invalid request received: %s\n", err)
	http.Error(w, err, http.StatusBadRequest)
}

func unprocessibleEntity(w http.ResponseWriter, err string) {
	log.Printf("ERR: Invalid request received: %s\n", err)
	http.Error(w, err, http.StatusUnprocessableEntity)
}

func notFound(w http.ResponseWriter) {
	log.Println("ERR: Resource not found.")
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func unsupportedMedia(w http.ResponseWriter) {
	err := fmt.Sprintf("%s. Expecting: %s: %s", http.StatusText(http.StatusUnsupportedMediaType), contentType, contentTypeJSON)
	log.Println(http.StatusText(http.StatusUnsupportedMediaType))
	http.Error(w, err, http.StatusUnsupportedMediaType)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %s\n", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
