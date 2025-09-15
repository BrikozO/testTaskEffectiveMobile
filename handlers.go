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

// CalculateSum godoc
//
//	@Summary		Calculate subscription sum
//	@Description	Calculate total sum for subscriptions in given period
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			calculation	body		dto.CalculationRequestDTO	true	"Calculation request"
//	@Success		200			{object}	object{price=integer}		"Returns total price"
//	@Failure		400				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/calculate [post]
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

// GetSubscriptions godoc
//
//	@Summary		Get user subscriptions
//	@Description	Get all subscriptions for a specific user
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string	true	"User ID (UUID)"	format(uuid)	example(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200		{array}		models.Subscription
//	@Failure		400				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/subscriptions/{user_id} [get]
func (app *application) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUserUuidFromRequest(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	subscriptions, err := app.subscriptions.GetByUserID(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// GetSubscriptionByID godoc
//
//	@Summary		Get subscription by ID
//	@Description	Get a specific subscription by user ID and subscription ID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string	true	"User ID (UUID)"	format(uuid)	example(550e8400-e29b-41d4-a716-446655440000)
//	@Param			subscription_id	path		int		true	"Subscription ID"
//	@Success		200				{object}	models.Subscription
//	@Failure		400				{string}	string
//	@Failure		404				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/subscriptions/{user_id}/{subscription_id} [get]
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

// PostSubscription godoc
//
//	@Summary		Create subscription
//	@Description	Create a new subscription. Dates should be in MM-YYYY format (e.g., "01-2024").
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		models.Subscription				true	"Subscription data"
//	@Success		201				{string}	string							"Created"
//	@Failure		400				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/subscriptions [post]
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

// UpdateSubscription godoc
//
//	@Summary		Update subscription
//	@Description	Update an existing subscription by ID. Dates should be in MM-YYYY format.
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription_id	path		int								true	"Subscription ID"
//	@Param			subscription	body		models.Subscription				true	"Updated subscription data"
//	@Success		202				{string}	string							"Accepted"
//	@Failure		400				{string}	string
//	@Failure		404				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/subscriptions/{subscription_id} [put]
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

// DeleteSubscription godoc
//
//	@Summary		Delete subscription
//	@Description	Delete a subscription by ID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription_id	path		int	true	"Subscription ID"
//	@Success		202				{object}	map[string]string
//	@Failure		400				{string}	string
//	@Failure		404				{string}	string
//	@Failure		500				{string}	string
//	@Router			/api/v1/subscriptions/{subscription_id} [delete]
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
