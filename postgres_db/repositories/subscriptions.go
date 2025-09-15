package repositories

import (
	"database/sql"
	"fmt"
	"testTaskEffectiveMobile/dto"
	"testTaskEffectiveMobile/models"

	"github.com/google/uuid"
)

type SubscriptionsRepository struct {
	Db *sql.DB
}

func (sr *SubscriptionsRepository) CalculateSum(calcDto dto.CalculationRequestDTO) (int64, error) {
	query := `
        SELECT COALESCE(SUM(price), 0) 
        FROM subscriptions 
        WHERE 1=1`

	var args []any
	argIndex := 1

	if calcDto.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *calcDto.UserID)
		argIndex++
	}

	if calcDto.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, *calcDto.ServiceName)
		argIndex++
	}

	query += fmt.Sprintf(` 
        AND start_date <= $%d 
        AND (end_date IS NULL OR end_date >= $%d)`,
		argIndex, argIndex+1)
	args = append(args, calcDto.EndDate, calcDto.StartDate)

	var totalCost int64
	err := sr.Db.QueryRow(query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return totalCost, nil
}

func (sr *SubscriptionsRepository) GetByUserID(userId uuid.UUID) ([]dto.SubscriptionDTO, error) {
	stmt := `SELECT id, service_name, price, user_id, start_date, end_date
			 FROM subscriptions
			 WHERE user_id = $1`

	rows, err := sr.Db.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []dto.SubscriptionDTO
	for rows.Next() {
		var s dto.SubscriptionDTO
		err = rows.Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (sr *SubscriptionsRepository) GetByUserIDAndID(userId uuid.UUID, id int) (dto.SubscriptionDTO, error) {
	stmt := `SELECT id, service_name, price, user_id, start_date, end_date
			 FROM subscriptions
			 WHERE user_id = $1
			 AND id = $2`

	var s dto.SubscriptionDTO

	err := sr.Db.QueryRow(stmt, userId, id).Scan(&s.Id, &s.ServiceName, &s.Price, &s.UserId, &s.StartDate, &s.EndDate)
	if err != nil {
		return dto.SubscriptionDTO{}, err
	}
	return s, nil
}

func (sr *SubscriptionsRepository) Insert(s models.Subscription) error {
	stmt := `INSERT INTO subscriptions(user_id, service_name, price, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5)`
	_, err := sr.Db.Exec(stmt, s.UserId, s.ServiceName, s.Price, s.StartDate, s.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SubscriptionsRepository) Update(id int, s models.Subscription) error {
	stmt := `update subscriptions
				set service_name = $2,
					user_id = $3,
					price = $4,
					start_date = $5,
					end_date = $6
				where id = $1`
	result, err := sr.Db.Exec(stmt, id, s.ServiceName, s.UserId, s.Price, s.StartDate, s.EndDate)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil

}

func (sr *SubscriptionsRepository) Delete(id int) error {
	result, err := sr.Db.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
