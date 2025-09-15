package dto

import (
	"testTaskEffectiveMobile/models"

	"github.com/google/uuid"
)

type CalculationRequestDTO struct {
	ServiceName *string              `json:"service_name"`
	UserID      *uuid.UUID           `json:"user_id"`
	StartDate   models.MonthYearDate `json:"start_date"`
	EndDate     models.MonthYearDate `json:"end_date"`
}
