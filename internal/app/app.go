package app

import (
	"crud/internal/config"
	"crud/internal/service"
	"log/slog"
	"net/http"
)

type App struct {
	Srv           service.Service
	Server        *http.Server
	HealthChecker service.Checker
	Config        *config.Config
	logger        slog.Logger
}

func New() *App {
	cfg := config.New()
	logger := slog.Default()

	return &App{
		Config: cfg,
		logger: *logger,
	}
}
