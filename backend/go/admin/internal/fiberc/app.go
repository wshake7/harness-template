package fiberc

import (
	"admin/internal/config"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	*fiber.App
	conf *config.Config
	done chan struct{}
}

func (app *App) Start() {
	go func() {
		err := app.Listen(fmt.Sprintf("%s:%d", app.conf.Host, app.conf.Port))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("Error starting server", zap.Error(err))
		}
	}()
	<-app.done
}
