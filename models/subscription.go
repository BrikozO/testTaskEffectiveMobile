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
	ServiceName string         `json:"service_name"`
	Price       int            `json:"price"`
	UserId      uuid.UUID      `json:"user_id"`
	StartDate   MonthYearDate  `json:"start_date"`
	EndDate     *MonthYearDate `json:"end_date"`
}
