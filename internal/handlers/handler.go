package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/suryansh74/sketlo/internal/config"
)

type GameHandler struct {
	cfg *config.Config
}

func NewGameHandler(cfg *config.Config) *GameHandler {
	return &GameHandler{
		cfg: cfg,
	}
}

func (gh *GameHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	view, err := gh.cfg.Views.GetTemplate("index.jet")
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}
	view.Execute(w, nil, nil)
}

func (gh *GameHandler) CheakGameHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "working fine",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (gh *GameHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	http.Redirect(w, r, "/game?username="+username, http.StatusSeeOther)
}
