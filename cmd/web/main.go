package main

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/suryansh74/sketlo/internal/config"
	"github.com/suryansh74/sketlo/internal/server"
)

func main() {
	views := jet.NewSet(
		jet.NewOSFileSystemLoader("./views"),
		jet.DevelopmentMode(true),
	)

	cfg := config.NewConfig(views)
	server := server.NewServer(cfg)
	server.StartServer()
}
