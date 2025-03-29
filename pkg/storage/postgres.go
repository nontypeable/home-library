package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"home-library/pkg/config"
)

func NewPostgres(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := cfg.GetDSN()
	return sqlx.Connect("postgres", dsn)
}
