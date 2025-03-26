package app

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"home-library/internal/server"
	"home-library/pkg/config"
	"home-library/pkg/storage"
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

	return &App{
		cfg:  cfg,
		echo: server.NewEchoServer(&cfg.HTTPServer),
		db:   db,
	}, nil
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
