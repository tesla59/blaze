package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/config"
)

// GetPool creates a new connection pool to the PostgreSQL database using the configuration from the config package.
func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := config.GetConfig()
	host := cfg.Db.Host
	port := cfg.Db.Port
	user := cfg.Db.User
	password := cfg.Db.Password
	database := cfg.Db.Dbname
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database)

	dbconf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbconf)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return pool, nil
}
