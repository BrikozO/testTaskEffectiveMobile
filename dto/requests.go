package dto

import (
	"testTaskEffectiveMobile/models"

	"github.com/google/uuid"
)

type CalculationRequestDTO struct {
	ServiceName *string              `json:"service_name" example:"Netflix"`
	UserID      *uuid.UUID           `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   models.MonthYearDate `json:"start_date" example:"01-2024" swaggertype:"string" format:"MM-YYYY"`
	EndDate     models.MonthYearDate `json:"end_date" example:"12-2024" swaggertype:"string" format:"MM-YYYY"`
}
