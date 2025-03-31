package testrepository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/simpleGorm/pg"
)

// Mandatory declare ID field
type TestObj struct {
	ID     int64
	Field1 int64
	Field2 string
}

type TestObjRepository struct {
	pg.Repository[TestObj]
}

const TABLE_NAME = "TestObjTable"

var increaseField1Builder = sq.Update(TABLE_NAME).PlaceholderFormat(sq.Dollar).
	Set("field1", sq.Expr("field1 + 1")).Suffix("RETURNING id, field1, field2")

func NewMyObjRepository(db pg.DbClient) TestObjRepository {
	return TestObjRepository{
		pg.NewPostgreRepository[TestObj](
			db,
			sq.Insert(TABLE_NAME).PlaceholderFormat(sq.Dollar).Columns("field1", "field2"),
			sq.Select("id", "field1", "field2").PlaceholderFormat(sq.Dollar).From(TABLE_NAME),
			sq.Update(TABLE_NAME).PlaceholderFormat(sq.Dollar),
			sq.Delete(TABLE_NAME).PlaceholderFormat(sq.Dollar),
			myObjConverter),
	}
}

func myObjConverter(row pgx.Row) TestObj {
	var myObj TestObj
	if err := row.Scan(&myObj.ID, &myObj.Field1, &myObj.Field2); err != nil {
		panic(err)
	}
	return myObj
}

func (repo *TestObjRepository) GetOneByField2(ctx context.Context, field2 string) TestObj {
	list := repo.GetBy(ctx, sq.Eq{"field2": field2})
	if list != nil && len(list) > 0 {
		return (list)[0]
	}
	var ret TestObj
	return ret
}

func (repo *TestObjRepository) IncreaseField1(ctx context.Context, id int64) int64 {
	updated := repo.UpdateReturning(ctx, increaseField1Builder.Where(sq.Eq{"id": id}))
	return updated.Field1
}
