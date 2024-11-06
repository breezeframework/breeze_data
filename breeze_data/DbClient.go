package breeze_data

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DbClient клиент для работы с БД
type DbClient interface {
	API() DbApi
	Close() error
}

// DbApi интерфейс для работы с БД
type DbApi interface {
	SQLExecutor
	Transactor
	Pinger
	Close()
}

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context) error

// TxManager менеджер транзакций, который выполняет указанный пользователем обработчик в транзакции
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться еще где-то, например, для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type SQLExecutor interface {
	NamedQueryExecutor
	QueryExecutor
}

// NamedQueryExecutor интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedQueryExecutor interface {
	/*ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error*/
}

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
