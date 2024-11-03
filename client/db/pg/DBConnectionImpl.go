package pg

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/breezeframework/breeze_data/client/db"
	"github.com/breezeframework/breeze_data/client/db/prettier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type key string

const (
	TxKey key = "tx"
)

type DBConnectionImpl struct {
	connectionPool *pgxpool.Pool
}

func NewDBConnection(pool *pgxpool.Pool) db.DBConnection {
	return &DBConnectionImpl{
		connectionPool: pool,
	}
}

/*func (p *DBConnectionImpl) ScanOneContext(ctx context.Context, dest base{}, q *sq.UserSelectBuilder) {
	row := p.QueryContextSelect(ctx, q, nil)
	err := pgxscan.ScanOne(dest, row)
	if err != nil {
		panic(err)
	}
}

func (p *DBConnectionImpl) ScanAllContext(ctx context.Context, dest base{}, q *sq.UserSelectBuilder) {
	rows := p.QueryContextSelect(ctx, q, nil)
	err := pgxscan.ScanAll(dest, rows)
	if err != nil {
		panic(err)
	}
}*/

func (p *DBConnectionImpl) ExecUpdate(ctx context.Context, builder *sq.UpdateBuilder) pgconn.CommandTag {
	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	log.Printf("[ExecUpdate] query: %s", query)
	log.Printf("[ExecUpdate] args: %+v", args)
	log.Printf("[ExecUpdate] err: %+v", err)
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	var tag pgconn.CommandTag
	if ok {
		tag, err = tx.Exec(ctx, query, args...)
	} else {
		tag, err = p.connectionPool.Exec(ctx, query, args...)
	}
	if err != nil {
		log.Printf("err: %+v", err)
		log.Panic(err)
	}
	return tag
}

func (p *DBConnectionImpl) QueryContextSelect(ctx context.Context, builder *sq.SelectBuilder, where map[string]interface{}) pgx.Rows {
	if where != nil {
		builder.Where(where)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated SQL query:", query)
	fmt.Println("Arguments:", args)
	fmt.Println("ctx:", ctx)
	fmt.Println("ctx.Value(TxKey):", ctx.Value(TxKey))
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	var rows pgx.Rows
	if ok {
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = p.connectionPool.Query(ctx, query, args...)
	}
	if err != nil {
		panic(err)
	}
	return rows
}

func (p *DBConnectionImpl) QueryRowContextSelect(ctx context.Context, builder *sq.SelectBuilder) pgx.Row {
	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated SQL query:", query)
	fmt.Println("Arguments:", args)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.connectionPool.QueryRow(ctx, query, args...)
}

func (p *DBConnectionImpl) QueryRowContextInsert(ctx context.Context, builder *sq.InsertBuilder) pgx.Row {

	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated SQL query:", query)
	fmt.Println("Arguments:", args)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}

	return p.connectionPool.QueryRow(ctx, query, args...)
}

func (p *DBConnectionImpl) BeginTx(ctx context.Context, txOptions pgx.TxOptions) pgx.Tx {
	tx, err := p.connectionPool.BeginTx(ctx, txOptions)
	if err != nil {
		panic(err)
	}
	return tx
}

func (p *DBConnectionImpl) Ping(ctx context.Context) error {
	return p.connectionPool.Ping(ctx)
}

func (p *DBConnectionImpl) Close() {
	p.connectionPool.Close()
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
