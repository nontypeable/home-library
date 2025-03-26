package main

import (
	"github.com/rs/zerolog/log"
	"home-library/internal/app"
	"home-library/pkg/config"
	"time"
	_ "time/tzdata"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load config")
	}

	_, err = time.LoadLocation(cfg.Application.TimeZone)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load timezone")
	}

	application, err := app.NewApp(*cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create application")
	}
	if err := application.Start(); err != nil {
		log.Fatal().Err(err).Msgf("failed to start application")
	}
}
