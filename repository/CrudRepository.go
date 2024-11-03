package repository

import (
	"context"
	"github.com/breezeframework/breeze_data/client/db"
)

type CrudRepository[T any] interface {
	GetDbConnection() db.DBConnection
	Create(ctx context.Context, entity T) (int64, error)
	GetById(ctx context.Context, id int64) (*T, error)
	GetAll(ctx context.Context) (*[]T, error)
	Update(ctx context.Context, id int64, entity T) error
	Delete(ctx context.Context, id int64) error
}
