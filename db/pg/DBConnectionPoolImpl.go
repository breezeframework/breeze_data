package pg

import (
	"context"
	"github.com/breezeframework/breeze_data/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgreDBConnectionPool struct {
	masterDBC db.DBConnection
}

func NewConnectionPool(ctx context.Context, dsn string) (db.DBConnectionPool, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &PostgreDBConnectionPool{
		masterDBC: &PostgreDBConnection{connectionPool: dbc},
	}, nil
}

func (c *PostgreDBConnectionPool) GetConnection() db.DBConnection {
	return c.masterDBC
}

func (c *PostgreDBConnectionPool) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
