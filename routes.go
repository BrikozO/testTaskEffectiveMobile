package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /calculate", app.calculateSum)
	router.HandleFunc("GET /subscriptions/{user_id}", app.getSubscriptions)
	router.HandleFunc("GET /subscriptions/{user_id}/{subscription_id}", app.getSubscriptionByID)
	router.HandleFunc("POST /subscriptions", app.postSubscription)
	router.HandleFunc("PUT /subscriptions/{subscription_id}", app.updateSubscription)
	router.HandleFunc("DELETE /subscriptions/{subscription_id}", app.deleteSubscription)
	return app.LogMiddleware(router)
}
