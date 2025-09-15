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

	_ "github.com/lib/pq"
)

// TODO: разобраться, зачем здесь указатели
type application struct {
	subscriptions *repositories.SubscriptionsRepository
	logger        *slog.Logger
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgres_db.Host,
		postgres_db.Port,
		postgres_db.User,
		postgres_db.Password,
		postgres_db.Db)
	db, err := postgres_db.ConnectPostgres(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
