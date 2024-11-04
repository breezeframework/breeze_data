package breeze_data

import (
	"context"
	sq "github.com/Masterminds/squirrel"
)

type CrudRepository[T any] interface {
	Create(ctx context.Context, entity T) (int64, error)
	GetById(ctx context.Context, id int64) (*T, error)
	GetAll(ctx context.Context) (*[]T, error)
	GetBy(ctx context.Context, where sq.Eq) (*[]T, error)
	Update(ctx context.Context, id int64, entity T) error
	Delete(ctx context.Context, id int64) error
}
