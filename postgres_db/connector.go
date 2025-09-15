package postgres_db

import (
	"database/sql"
	"errors"
	"log/slog"
)

func ConnectPostgres(connectionString string) (*sql.DB, func(), error) {
	var err error
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, nil, errors.New("could not connect to postgres: " + err.Error())
	}
	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, nil, errors.New("could not connect to postgres: " + err.Error())
	}
	slog.Info("connected to postgres")
	return conn, func() {
		conn.Close()
	}, nil
}
