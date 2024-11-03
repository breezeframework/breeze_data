package db

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBConnection интерфейс для работы с БД
type DBConnection interface {
	SQLExecutor
	Transactor
	Pinger
	Close()
}

// Query обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться еще где-то, например, для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecutor комбинирует NamedQueryExecutor и QueryExecutor
type SQLExecutor interface {
	NamedQueryExecutor
	QueryExecutor
}

// NamedQueryExecutor интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedQueryExecutor interface {
	/*	ScanOneContext(ctx context.Context, dest base{}, q *squirrel.UserSelectBuilder)
		ScanAllContext(ctx context.Context, dest base{}, q *squirrel.UserSelectBuilder)*/
}

// QueryExecutor интерфейс для работы с обычными запросами
type QueryExecutor interface {
	ExecUpdate(ctx context.Context, builder *squirrel.UpdateBuilder) pgconn.CommandTag
	QueryContextSelect(ctx context.Context, builder *squirrel.SelectBuilder, where map[string]interface{}) pgx.Rows
	QueryRowContextSelect(ctx context.Context, builder *squirrel.SelectBuilder) pgx.Row
	QueryRowContextInsert(ctx context.Context, builder *squirrel.InsertBuilder) pgx.Row
}

// Pinger интерфейс для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}
