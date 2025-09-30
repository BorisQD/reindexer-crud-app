package app

import (
	api "crud/internal/api/http"
	"crud/internal/client"
	"crud/internal/service"
	"net"
	"net/http"
)

func (app *App) Bootstrap() {
	db := client.New(app.Config.DB)
	srv := service.New(db, app.Config.TTL)
	checker := service.NewChecker(db)

	app.Srv = *srv
	app.HealthChecker = *checker

	server := api.NewStrictHandler(api.Server{
		Service: app.Srv,
		Checker: app.HealthChecker,
		Logger:  app.logger.With("api", "http"),
	}, nil)

	r := http.NewServeMux()

	h := api.HandlerFromMux(server, r)
	app.Server = &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort(app.Config.Server.Host, app.Config.Server.Port),
	}
}
