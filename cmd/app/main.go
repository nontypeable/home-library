package main

import (
	"github.com/rs/zerolog/log"
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
}
