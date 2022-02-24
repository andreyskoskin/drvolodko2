package webapi

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"github.com/andreyskoskin/drvolodko2/audit"
)

type (
	Config struct {
		HTTP HTTPConfig `toml:"http"`
		DB   DBConfig   `toml:"db"`
	}

	HTTPConfig struct {
		Address string `toml:"address"`
	}

	DBConfig struct {
		Name     string `toml:"name"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
	}
)

func Start(c Config) error {
	var db, err = sql.Open("postgres", "<connection string>")
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	var (
		pg = &postgres{db}
		e  = echo.New()
	)

	newAuditAPI(pg).install(e.Group("/audit"))
	return e.Start(c.HTTP.Address)
}

type postgres struct {
	db *sql.DB
}

func (pg *postgres) Exec(query string, params ...interface{}) (sql.Result, error) {
	return pg.db.Exec(query, params...)
}

func (pg *postgres) WithTx(ctx context.Context, fn func(audit.SqlDB) error, opts *sql.TxOptions) (err error) {
	var tx *sql.Tx

	tx, err = pg.db.BeginTx(ctx, opts)
	if err != nil {
		return
	}

	// nolint
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
