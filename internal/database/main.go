package database

import (
	"fmt"
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

func Connect(username string, password string, host string, table string) (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, host, table))

	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(stdlib.OpenDB(*connConfig), "pgx"), nil
}
