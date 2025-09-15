package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /calculate", app.calculateSum)
	router.HandleFunc("GET /subscriptions/{user_id}", app.getSubscriptions)
	router.HandleFunc("GET /subscriptions/{user_id}/{subscription_id}", app.getSubscriptionByID)
	router.HandleFunc("POST /subscriptions", app.postSubscription)
	router.HandleFunc("PUT /subscriptions/{subscription_id}", app.updateSubscription)
	router.HandleFunc("DELETE /subscriptions/{subscription_id}", app.deleteSubscription)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", app.LogMiddleware(router)))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return mux
}
