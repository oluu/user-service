package main

import (
	"net/http"

	"github.com/shinypotato/user-service/data"
	"github.com/shinypotato/user-service/handlers"
	"github.com/shinypotato/user-service/service"
)

const port = ":3000"

func main() {
	repository := data.InitRepository()
	userService := service.NewUserService(repository)
	handlers.RegisterHandlers(userService)
	http.ListenAndServe(port, nil)
}
