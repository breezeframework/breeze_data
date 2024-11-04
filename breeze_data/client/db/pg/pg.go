package pg

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/breezeframework/breeze_data/breeze_data/client/db"
	"github.com/breezeframework/breeze_data/breeze_data/client/db/prettier"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type key string

const (
	TxKey key = "tx"
)

type pg struct {
	dbc *pgxpool.Pool
}

func (p *pg) ExecUpdate(ctx context.Context, builder *sq.UpdateBuilder) pgconn.CommandTag {
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
		tag, err = p.dbc.Exec(ctx, query, args...)
	}
	if err != nil {
		log.Printf("err: %+v", err)
		log.Panic(err)
	}
	return tag
}

func (p *pg) QueryContextSelect(ctx context.Context, builder *sq.SelectBuilder, where map[string]interface{}) pgx.Rows {
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
		rows, err = p.dbc.Query(ctx, query, args...)
	}
	if err != nil {
		panic(err)
	}
	return rows
}

func (p *pg) QueryRowContextSelect(ctx context.Context, builder *sq.SelectBuilder) pgx.Row {
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

	return p.dbc.QueryRow(ctx, query, args...)
}

func (p *pg) QueryRowContextInsert(ctx context.Context, builder *sq.InsertBuilder) pgx.Row {

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

	return p.dbc.QueryRow(ctx, query, args...)
}

func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

/*func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}*/

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
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
