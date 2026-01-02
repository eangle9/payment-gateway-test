package persistencedb

import (
	"context"
	"pg/internal/constant/model/db"
	"pg/platform/hlog"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PersistenceDB struct {
	*db.Queries
	pool    *pgxpool.Pool
	log     hlog.Logger
	options Options
}

type Options struct {
	SSODB        Sibling
	AuthzDB      Sibling
	AccountingDB Sibling
}

func setOptions(options Options) Options {
	if len(options.SSODB) == 0 {
		options.SSODB = "sso"
	}

	if len(options.AuthzDB) == 0 {
		options.AuthzDB = "authz"
	}

	if len(options.AccountingDB) == 0 {
		options.AccountingDB = "accounting"
	}

	return options
}

type Sibling string

func (q PersistenceDB) WithTransaction(ctx context.Context, fn func(tx PersistenceDB) error) error {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(q.WithTx(tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (q PersistenceDB) WithTx(tx db.DBTX) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(tx),
		pool:    q.pool,
		log:     q.log,
		options: q.options,
	}
}

func New(pool *pgxpool.Pool, log hlog.Logger, options Options) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		pool:    pool,
		log:     log,
		options: setOptions(options),
	}
}
