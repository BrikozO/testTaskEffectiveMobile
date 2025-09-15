package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var monthYearDateFormat = "01-2006"

// MonthYearDate represents a date in MM-YYYY format
// swagger:strfmt month-year
type MonthYearDate struct {
	time.Time
}

func (m *MonthYearDate) MarshalJSON() ([]byte, error) {
	if m == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(m.Format(monthYearDateFormat))
}

func (m *MonthYearDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return nil
	}
	t, err := time.Parse(monthYearDateFormat, s)
	if err != nil {
		return err
	}
	*m = MonthYearDate{t}
	return nil
}

// TODO: разобраться подробнее с работой этих ресиверов
func (m MonthYearDate) Value() (driver.Value, error) {
	return m.Time, nil
}

func (m *MonthYearDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*m = MonthYearDate{v}
		return nil
	case string:
		t, err := time.Parse(monthYearDateFormat, v)
		if err != nil {
			return err
		}
		*m = MonthYearDate{t}
		return nil
	default:
		return errors.New("incompatible type for MonthYearDate")
	}
}

type Subscription struct {
	ServiceName string         `json:"service_name" example:"Netflix"`
	Price       int            `json:"price" example:"999"`
	UserId      uuid.UUID      `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   MonthYearDate  `json:"start_date" example:"01-2024" swaggertype:"string" format:"MM-YYYY"`
	EndDate     *MonthYearDate `json:"end_date" example:"12-2024" swaggertype:"string" format:"MM-YYYY"`
}
