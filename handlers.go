package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testTaskEffectiveMobile/dto"
	"testTaskEffectiveMobile/models"

	"github.com/google/uuid"
)

func parseUserUuidFromRequest(r *http.Request) (uuid.UUID, error) {
	userId := r.PathValue("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}
	return uid, nil
}

func (app *application) calculateSum(w http.ResponseWriter, r *http.Request) {
	var calcDto dto.CalculationRequestDTO
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&calcDto)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	if calcDto.StartDate.IsZero() || calcDto.EndDate.IsZero() {
		http.Error(w, "Period start and end required", http.StatusBadRequest)
		return
	}
	result, err := app.subscriptions.CalculateSum(calcDto)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"price": "%d"}`, result)))
}

func (app *application) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUserUuidFromRequest(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	subscriptions, err := app.subscriptions.GetByUserID(userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

func (app *application) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUserUuidFromRequest(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	subscriptionId := r.PathValue("subscription_id")
	intSubscrId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	s, err := app.subscriptions.GetByUserIDAndID(userId, intSubscrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (app *application) postSubscription(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sub models.Subscription
	err := decoder.Decode(&sub)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.subscriptions.Insert(sub)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *application) updateSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionId := r.PathValue("subscription_id")
	intSubscrId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var sub models.Subscription
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&sub)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.subscriptions.Update(intSubscrId, sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionId := r.PathValue("subscription_id")
	intSubscrId, err := strconv.Atoi(subscriptionId)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err = app.subscriptions.Delete(intSubscrId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"detail": "%s"}`, "subscription successfully deleted")))
	return
}
