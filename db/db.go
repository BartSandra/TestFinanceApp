package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

// InitDB initializes a database connection.
func InitDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	return db, nil
}
