package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"testTaskEffectiveMobile/postgres_db"
	"testTaskEffectiveMobile/postgres_db/repositories"
	"time"

	_ "testTaskEffectiveMobile/docs"

	_ "github.com/lib/pq"
)

// TODO: разобраться, зачем здесь указатели
type application struct {
	subscriptions *repositories.SubscriptionsRepository
	logger        *slog.Logger
}

//	@title			Swagger API Documentation
//	@version		1.0.0
//	@description	Swagger for test Task in effective Mobile

//	@contact.name	Oleg (API Author):
//	@contact.url	https://github.com/BrikozO
//	@contact.email	oleg.yakushev.work@gmail.com

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db, closer, err := postgres_db.ConnectPostgres(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	app := &application{subscriptions: &repositories.SubscriptionsRepository{Db: db},
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil))}

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      app.routes(),
	}
	err = s.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
