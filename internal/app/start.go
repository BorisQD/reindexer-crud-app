package app

import (
	"context"
	"fmt"
)

func Start(ctx context.Context) error {
	app := New()
	err := app.Config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	app.Bootstrap()

	app.logger.Info("starting app")

	if err := app.Srv.Start(ctx); err != nil {
		return fmt.Errorf("start service: %w", err)
	}

	defer app.Shutdown()

	return app.Server.ListenAndServe()
}
