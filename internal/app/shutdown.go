package app

import (
	"context"
	"fmt"
	"time"
)

const shutdownDuration = 5 * time.Second

func (app *App) Shutdown() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	if app.server != nil {
		if err := app.server.Shutdown(ctxWithTimeout); err != nil {
			app.logger.Error(fmt.Errorf("Server.Shutdown: %w", err).Error())
		}
	}

	app.srv.Close(ctxWithTimeout)

	<-ctxWithTimeout.Done()
	app.logger.Info(fmt.Sprintf("Shutdown timeout of %v", shutdownDuration))
	cancel()
}
