package database

import (
	"TelegramGPT/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Executor interface {
	sqlx.ExtContext
	sqlx.Ext

	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

func Connect(config config.DatabaseConfig) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(config.Url)

	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(stdlib.OpenDB(*connConfig), "pgx"), nil
}
