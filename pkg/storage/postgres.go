package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"home-library/pkg/config"
)

func NewPostgres(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
