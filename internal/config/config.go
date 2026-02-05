package config

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/suryansh74/sketlo/internal/chat"
)

type Config struct {
	Views *jet.Set
	Hub   *chat.Hub
}

func NewConfig(views *jet.Set, hub *chat.Hub) *Config {
	return &Config{
		Views: views,
		Hub:   hub,
	}
}
