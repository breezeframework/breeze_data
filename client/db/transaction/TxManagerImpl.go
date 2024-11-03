package transaction

import (
	"context"
	"github.com/breezeframework/breeze_data/client/db"
	"github.com/breezeframework/breeze_data/client/db/pg"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type TxManagerImpl struct {
	db  db.Transactor
	ctx context.Context
}

// NewTransactionManager создает новый менеджер транзакций, который удовлетворяет интерфейсу db.TxManager
func NewTransactionManager(ctx context.Context, db db.Transactor) db.TxManager {
	return &TxManagerImpl{
		db:  db,
		ctx: ctx,
	}
}

// transaction основная функция, которая выполняет указанный пользователем обработчик в транзакции
func (m *TxManagerImpl) transaction(opts pgx.TxOptions, fn db.Handler) {
	// Если это вложенная транзакция, пропускаем инициацию новой транзакции и выполняем обработчик.
	tx, ok := m.ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		fn(m.ctx)
	}

	// Стартуем новую транзакцию.
	tx = m.db.BeginTx(m.ctx, opts)

	// Кладем транзакцию в контекст.
	ctx := pg.MakeContextTx(m.ctx, tx)

	// Настраиваем функцию отсрочки для отката или коммита транзакции.
	defer func() {
		var err error
		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// откатываем транзакцию, если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		// если ошибок не было, коммитим транзакцию
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Выполните код внутри транзакции.
	// Если функция терпит неудачу, возвращаем ошибку, и функция отсрочки выполняет откат
	// или в противном случае транзакция коммитится.
	fn(ctx)

}

func (m *TxManagerImpl) ReadCommitted(f db.Handler) {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadUncommitted}
	m.transaction(txOpts, f)
}
