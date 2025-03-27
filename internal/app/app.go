package app

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"home-library/internal/server"
	"home-library/internal/services/user/entities"
	"home-library/pkg/config"
	"home-library/pkg/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	db   *gorm.DB
	echo *echo.Echo
	cfg  config.Config
}

func NewApp(cfg config.Config) (*App, error) {
	db, err := storage.NewPostgres(&cfg.Database)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:  cfg,
		echo: server.NewEchoServer(&cfg.HTTPServer),
		db:   db,
	}, nil
}

func (app *App) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.cfg.Application.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := app.echo.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
		} else {
			log.Info().Msg("server was successfully shutdown")
		}

		err := app.closeDatabaseConnection()
		if err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		} else {
			log.Info().Msg("database connection was successfully closed")
		}
	}()

	if err := app.startService(); err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", app.cfg.HTTPServer.Host, app.cfg.HTTPServer.Port)
	if err := app.echo.StartTLS(address, app.cfg.SSL.CertFile, app.cfg.SSL.KeyFile); err != nil {
		return err
	}

	return nil
}

func (app *App) closeDatabaseConnection() error {
	db, err := app.db.DB()
	if err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
