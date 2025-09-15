package postgres_db

import (
	"database/sql"
	"errors"
	"log/slog"
)

const (
	Host     = "localhost"
	Port     = 5432
	User     = "testuser"
	Password = "12345678"
	Db       = "test_db"
)

func ConnectPostgres(connectionString string) (*sql.DB, error) {
	var err error
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, errors.New("could not connect to postgres: " + err.Error())
	}
	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, errors.New("could not connect to postgres: " + err.Error())
	}
	slog.Info("connected to postgres")
	return conn, nil
}
