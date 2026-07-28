package postgres

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
)




func ConnectionDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {

	dbPool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return  nil,err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		dbPool.Close()
		return nil, err
	}

	return dbPool, nil
}