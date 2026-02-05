package server

import (
	"github.com/suryansh74/sketlo/internal/handlers"
)

func (srv *server) SetupRoutes(gameHandler *handlers.GameHandler) {
	// web endpoints from here
	srv.router.Get("/", gameHandler.HomeHandler)
	srv.router.Post("/join", gameHandler.JoinRoom)
	srv.router.Get("/game", gameHandler.GameRoom)

	// api endpoints from here
	srv.router.Get("/api/check_health", gameHandler.CheakGameHealth)
}
