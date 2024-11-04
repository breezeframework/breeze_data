package db

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context) error

// Client клиент для работы с БД
type Client interface {
	DB() DB
	Close() error
}

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

// Transactor интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedExecer interface {
	/*ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error*/
}

type QueryExecer interface {
	ExecUpdate(ctx context.Context, builder *squirrel.UpdateBuilder) pgconn.CommandTag
	QueryContextSelect(ctx context.Context, builder *squirrel.SelectBuilder, where map[string]interface{}) pgx.Rows
	QueryRowContextSelect(ctx context.Context, builder *squirrel.SelectBuilder) pgx.Row
	QueryRowContextInsert(ctx context.Context, builder *squirrel.InsertBuilder) pgx.Row
}

// Pinger интерфейс для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB интерфейс для работы с БД
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
