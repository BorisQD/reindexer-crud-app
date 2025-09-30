package app

import (
	"context"
	api "crud/internal/api/http"
	"crud/internal/client"
	"crud/internal/domain"
	"crud/internal/service"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

func Start(ctx context.Context) error {
	db := client.New()
	s := service.New(db)

	app := New()
	app.logger = *slog.Default()
	app.srv = *s

	server := api.NewStrictHandler(api.Server{}, nil)

	r := http.NewServeMux()

	h := api.HandlerFromMux(server, r)
	app.server = &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort("0.0.0.0", "3000"),
	}

	defer app.Shutdown()

	app.logger.Info("starting app")

	if err := app.srv.Start(ctx); err != nil {
		return fmt.Errorf("start service: %w", err)
	}
	//return T(app)
	return app.server.ListenAndServe()
}

func T(a *App) error {
	if err := a.srv.CreateItem(context.Background(), domain.Item{Name: "First"}); err != nil {
		return err
	}
	if items, totalCount, err := a.srv.GetItemsPaginated(context.Background(), domain.Pagination{Offset: 0, Limit: 100}); err != nil {
		return err
	} else {
		a.logger.Debug(fmt.Sprintf("count: %d, first name: %s", totalCount, items[0].Name))
	}

	return nil
}
