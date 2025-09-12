package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type monthYearDate struct {
	time.Time
}

func (m *monthYearDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return errors.New("monthYearDate: cannot unmarshal null")
	}
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	*m = monthYearDate{t}
	return nil
}

func parseUuidFromRequest(r *http.Request) (uuid.UUID, error) {
	userId := r.PathValue("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}

type Subscription struct {
	ServiceName string         `json:"service_name"`
	Price       int            `json:"price"`
	ID          uuid.UUID      `json:"user_id"`
	StartDate   monthYearDate  `json:"start_date"`
	EndDate     *monthYearDate `json:"end_date"`
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get subscribes req"))
}

func getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUuidFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, err.Error())))
		return
	}
	row := conn.QueryRow(`SELECT service_name, price, user_id, start_date, end_date
FROM subscriptions
WHERE user_id = $1`, userId)
	var s Subscription
	var startDate time.Time
	var endDate sql.NullTime
	if err := row.Scan(&s.ServiceName, &s.Price, &s.ID, &startDate, &endDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, "Subscription not found")))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, err.Error())))
		return
	}
	s.StartDate = monthYearDate{startDate}
	if endDate.Valid {
		s.EndDate = &monthYearDate{endDate.Time}
	}
	json.NewEncoder(w).Encode(s)
}

func postSubscription(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sub Subscription
	err := decoder.Decode(&sub)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var endDate sql.NullTime
	if sub.EndDate != nil {
		endDate = sql.NullTime{Time: sub.EndDate.Time, Valid: true}
	} else {
		endDate = sql.NullTime{}
	}
	_, err = conn.Exec(`INSERT INTO subscriptions(user_id, service_name, price, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5)`, sub.ID, sub.ServiceName, sub.Price, sql.NullTime{Time: sub.StartDate.Time, Valid: true}, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func updateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update subscribes req"))
}

func deleteSubscription(w http.ResponseWriter, r *http.Request) {
	uid, err := parseUuidFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, err.Error())))
		return
	}
	res, err := conn.Exec("DELETE FROM subscriptions WHERE user_id=$1", uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, err.Error())))
		return
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, err.Error())))
		return
	}
	if deleted == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, "subscription not found")))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, "subscription successfully deleted")))
	return
}

const (
	host     = "localhost"
	port     = 5432
	user     = "testuser"
	password = "12345678"
	db       = "test_db"
)

var conn *sql.DB

func connectPostgres(connectionString string) error {
	var err error
	conn, err = sql.Open("postgres", connectionString)
	if err != nil {
		return errors.New("could not connect to postgres: " + err.Error())
	}
	defer conn.Close()
	err = conn.Ping()
	if err != nil {
		return errors.New("could not connect to postgres: " + err.Error())
	}
	slog.Info("connected to postgres")
	return nil
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db)
	err := connectPostgres(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	router.HandleFunc("GET /subscriptions", getSubscriptions)
	router.HandleFunc("GET /subscriptions/{user_id}", getSubscriptionByID)
	router.HandleFunc("POST /subscriptions", postSubscription)
	router.HandleFunc("PUT /subscriptions", updateSubscription)
	router.HandleFunc("DELETE /subscriptions/{user_id}", deleteSubscription)

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}
	err = s.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}
