package v1

import (
	"database/sql"
	"log/slog"
)

type SubscriptionMigration struct {
	Db *sql.DB
}

func (sb *SubscriptionMigration) Init() error {
	stmt := `create table if not exists subscriptions
(
    id           serial
        primary key,
    service_name varchar(256)             not null,
    price        integer                  not null,
    user_id      uuid                     not null,
    start_date   timestamp with time zone not null,
    end_date     timestamp with time zone
);`
	_, err := sb.Db.Exec(stmt)
	if err != nil {
		return err
	}
	slog.Info("Subscription migration v1 initialized")
	return nil
}
