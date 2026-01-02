package persistencedb

import (
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

func New(pool *pgxpool.Pool, log hlog.Logger, options Options) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		pool:    pool,
		log:     log,
		options: setOptions(options),
	}
}
