package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/suryansh74/sketlo/internal/chat"
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

// Game Route
// ==================================================

// HomeHandler get home page
func (gh *GameHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	view, err := gh.cfg.Views.GetTemplate("index.jet")
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}
	view.Execute(w, nil, nil)
}

// JoinRoom post with username
func (gh *GameHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	http.Redirect(w, r, "/game?username="+username, http.StatusSeeOther)
}

func (gh *GameHandler) GameRoom(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	fmt.Println(username)

	view, err := gh.cfg.Views.GetTemplate("game.jet")
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
	}
	getConnectedUsers := gh.cfg.Hub.GetConnectedUsers()
	vars := make(jet.VarMap)
	vars.Set("Username", username)
	vars.Set("ConnectedUsers", getConnectedUsers)

	view.Execute(w, vars, nil)
}

func (gh *GameHandler) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// use client function
	chat.ServeWs(gh.cfg.Hub, w, r)
}

// Api Route
// ==================================================

// CheakGameHealth api check
func (gh *GameHandler) CheakGameHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "working fine",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Helper function
// ==================================================
