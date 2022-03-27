package connections

import (
	"context"

	"emperror.dev/errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func GetPostgresDatabase(ctx context.Context, connectionString string) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(ctx, "pgx", connectionString)
	if err != nil {
		err = errors.WithMessage(err, "connecting")
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		err = errors.WithMessage(err, "ping")
		return nil, err
	}

	return db, nil

}
