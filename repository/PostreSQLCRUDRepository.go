package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/breezeframework/breeze_data/client/db"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const (
	idColumn     = "id"
	RETURNING_ID = "RETURNING id"
)

type PostgreSQLCRUDRepository[T any] struct {
	dbConnection  db.DBConnection
	insertBuilder sq.InsertBuilder
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
	deleteBuilder sq.DeleteBuilder
	scanner       func(row pgx.Row) (*T, error)
}

func (repo *PostgreSQLCRUDRepository[T]) GetDbConnection() db.DBConnection {
	return repo.dbConnection
}

func NewPostgreSQLCRUDRepository[T any](
	dbConnection db.DBConnection,
	insertBuilder sq.InsertBuilder,
	selectBuilder sq.SelectBuilder,
	updateBuilder sq.UpdateBuilder,
	deleteBuilder sq.DeleteBuilder,
	scanner func(pgx.Row) (*T, error)) CrudRepository[T] {
	return &PostgreSQLCRUDRepository[T]{
		dbConnection:  dbConnection,
		insertBuilder: insertBuilder, selectBuilder: selectBuilder, updateBuilder: updateBuilder, deleteBuilder: deleteBuilder,
		scanner: scanner}
}

func (repo *PostgreSQLCRUDRepository[T]) Create(ctx context.Context, entity T) (int64, error) {
	builder := repo.insertBuilder.Suffix(RETURNING_ID).Values(entity)
	var id int64
	err := repo.dbConnection.QueryRowContextInsert(ctx, &builder).Scan(&id)
	return id, err
}

func (repo *PostgreSQLCRUDRepository[T]) GetById(ctx context.Context, id int64) (*T, error) {
	builder := repo.selectBuilder.Where(sq.Eq{idColumn: id})
	row := repo.dbConnection.QueryRowContextSelect(ctx, &builder)
	return repo.scanner(row)
}

func (repo *PostgreSQLCRUDRepository[T]) ConvertToObjects(rows pgx.Rows) (*[]T, error) {
	var objs []T
	for rows.Next() {
		obj, err := repo.scanner(rows)
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
	rows := repo.dbConnection.QueryContextSelect(ctx, &repo.selectBuilder, nil)
	objs, err := repo.ConvertToObjects(rows)
	return objs, err
}

func (repo *PostgreSQLCRUDRepository[T]) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (repo *PostgreSQLCRUDRepository[T]) Update(ctx context.Context, id int64, entity T) error {
	panic("implement me")
}
