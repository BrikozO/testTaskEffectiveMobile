package dto

import "testTaskEffectiveMobile/models"

type SubscriptionDTO struct {
	Id int `json:"id"`
	models.Subscription
}
