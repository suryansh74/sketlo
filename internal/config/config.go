package config

import "github.com/CloudyKit/jet/v6"

type Config struct {
	Views *jet.Set
}

func NewConfig(views *jet.Set) *Config {
	return &Config{
		Views: views,
	}
}
