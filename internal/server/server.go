package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/suryansh74/sketlo/internal/config"
	"github.com/suryansh74/sketlo/internal/handlers"
)

type server struct {
	router *chi.Mux
	cfg    *config.Config
}

func NewServer(cfg *config.Config) *server {
	server := &server{
		router: chi.NewRouter(),
		cfg:    cfg,
	}
	server.router.Use(middleware.Logger)

	gameHandler := handlers.NewGameHandler(server.cfg)
	server.SetupRoutes(gameHandler)
	return server
}

func (srv *server) StartServer() {
	http.ListenAndServe("localhost:8000", srv.router)
}
