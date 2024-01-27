package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	QueryExecer
	TxStarter
	Close()
}

type QueryExecer interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}
type TxStarter interface {
	BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error)
}

func ExecTX(ctx context.Context, db TxStarter, txHandler func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()
	err = txHandler(ctx, tx)

	return err
}

type Tx interface {
	pgx.Tx
}

type Row interface {
	pgx.Row
}

type Rows interface {
	pgx.Rows
}
