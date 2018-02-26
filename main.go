package main

import (
	"net/http"

	"github.com/oluu/user-service/data"
	"github.com/oluu/user-service/handlers"
	"github.com/oluu/user-service/message"
	"github.com/oluu/user-service/service"
)

const port = ":3000"

func main() {
	repository := data.InitRepository()
	producer, consumer := message.InitMessaging()
	userService := service.NewUserService(repository, producer)
	message.InitHandlers(consumer, repository)
	handlers.RegisterHandlers(userService)
	http.ListenAndServe(port, nil)
}
