package pg

import (
	"context"
	"github.com/breezeframework/breeze_data/breeze_data/client/db"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkg/errors"
)

type pgDbClient struct {
	masterDBC db.DbApi
}

func NewPgDBClient(ctx context.Context, dsn string) (db.DbClient, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgDbClient{
		masterDBC: &pg{api: dbc},
	}, nil
}

func (c *pgDbClient) API() db.DbApi {
	return c.masterDBC
}

func (c *pgDbClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
