package server

import (
	"github.com/loanengine/internal/handler/rest"
)

type App struct {
	*Config
	RestHandler rest.RestHandler
}

type Config struct {
	Host string
	Port int
}

func NewApp(cfg *Config) *App {
	return &App{
		Config: cfg,
	}
}

func NewConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
}
