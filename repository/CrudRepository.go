package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/breezeframework/breeze_data/client/db"
)

type CrudRepository[T any] interface {
	GetDbConnection() db.DBConnection
	Create(ctx context.Context, entity T) (int64, error)
	GetById(ctx context.Context, id int64) (*T, error)
	GetAll(ctx context.Context) (*[]T, error)
	GetWhere(ctx context.Context, where sq.Eq) (*[]T, error)
	Update(ctx context.Context, id int64, entity T) error
	Delete(ctx context.Context, id int64) error
}
