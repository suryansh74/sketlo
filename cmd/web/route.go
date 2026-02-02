package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// api endpoints from here
	r.Get("/api/check_health", CheakHealth)
	return r
}

func CheakHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "working fine",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
