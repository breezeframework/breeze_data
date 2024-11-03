package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context)

// TxManager менеджер транзакций, который выполняет указанный пользователем обработчик в транзакции
type TxManager interface {
	ReadCommitted(f Handler)
}

// Transactor интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) pgx.Tx
}
