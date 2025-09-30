package app

import (
	"crud/internal/service"
	"log/slog"
	"net/http"
)

type App struct {
	srv    service.Service
	server *http.Server
	logger slog.Logger
}

func New() *App {
	return &App{}
}
