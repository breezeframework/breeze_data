package breeze_data

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	idColumn     = "id"
	RETURNING_ID = "RETURNING id"
)

type PostgreSQLCRUDRepository[T any] struct {
	db              DbClient
	insertBuilder   sq.InsertBuilder
	selectBuilder   sq.SelectBuilder
	updateBuilder   sq.UpdateBuilder
	deleteBuilder   sq.DeleteBuilder
	entityConverter func(row pgx.Row) (*T, error)
}

func NewPostgreSQLCRUDRepository[T any](
	db DbClient,
	insertBuilder sq.InsertBuilder,
	selectBuilder sq.SelectBuilder,
	updateBuilder sq.UpdateBuilder,
	deleteBuilder sq.DeleteBuilder,
	entityConverter func(pgx.Row) (*T, error)) CrudRepository[T] {
	return &PostgreSQLCRUDRepository[T]{
		db:            db,
		insertBuilder: insertBuilder, selectBuilder: selectBuilder, updateBuilder: updateBuilder, deleteBuilder: deleteBuilder,
		entityConverter: entityConverter}
}

func (repo *PostgreSQLCRUDRepository[T]) Create(ctx context.Context, entity T) (int64, error) {
	builder := repo.insertBuilder.Suffix(RETURNING_ID).Values(entity)
	var id int64
	err := repo.db.API().QueryRowContextInsert(ctx, &builder).Scan(&id)
	return id, err
}

func (repo *PostgreSQLCRUDRepository[T]) GetById(ctx context.Context, id int64) (*T, error) {
	builder := repo.selectBuilder.Where(sq.Eq{idColumn: id})
	row := repo.db.API().QueryRowContextSelect(ctx, &builder)
	return repo.entityConverter(row)
}

func (repo *PostgreSQLCRUDRepository[T]) ConvertToObjects(rows pgx.Rows) (*[]T, error) {
	var objs []T
	for rows.Next() {
		obj, err := repo.entityConverter(rows)
		if err != nil {
			return nil, err
		}
		objs = append(objs, *obj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &objs, nil
}

func (repo *PostgreSQLCRUDRepository[T]) GetAll(ctx context.Context) (*[]T, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	rows := repo.db.API().QueryContextSelect(ctx, &repo.selectBuilder, nil)
	objs, err := repo.ConvertToObjects(rows)
	return objs, err
}

func (repo *PostgreSQLCRUDRepository[T]) GetBy(ctx context.Context, where sq.Eq) (*[]T, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()
	builder := repo.selectBuilder.Where(where)
	rows := repo.db.API().QueryContextSelect(ctx, &builder, nil)
	objs, err := repo.ConvertToObjects(rows)
	return objs, err
}

func (repo *PostgreSQLCRUDRepository[T]) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (repo *PostgreSQLCRUDRepository[T]) Update(ctx context.Context, id int64, entity T) error {
	panic("implement me")
}
