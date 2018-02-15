package main

import (
	"net/http"

	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/handlers"
	"github.com/shinypotato/user-service/message"
	"github.com/shinypotato/user-service/service"
)

const port = ":3000"

func main() {
	repository := data.InitRepository()
	producer, consumers := message.InitMessaging([]string{"user-create", "user-update", "user-delete"})
	userService := service.NewUserService(repository, producer)
	message.InitHandlers(consumers, repository)
	handlers.RegisterHandlers(userService)
	http.ListenAndServe(port, nil)
}
