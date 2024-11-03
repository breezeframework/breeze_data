package pg

import (
	"context"
	"github.com/breezeframework/breeze_data/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DBConnectionPoolImpl struct {
	masterDBC db.DBConnection
}

func NewConnectionPool(ctx context.Context, dsn string) (db.DBConnectionPool, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &DBConnectionPoolImpl{
		masterDBC: &DBConnectionImpl{connectionPool: dbc},
	}, nil
}

func (c *DBConnectionPoolImpl) GetConnection() db.DBConnection {
	return c.masterDBC
}

func (c *DBConnectionPoolImpl) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
