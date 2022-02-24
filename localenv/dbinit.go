package localenv

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/andreyskoskin/drvolodko2/webapi"
)

var ddls = []string{
	`CREATE SCHEMA IF NOT EXISTS audit`,
	`CREATE TABLE IF NOT EXISTS audit.audit (
		id       INT    NOT NULL,
		num      INT    NOT NULL,
		month    INT    NOT NULL,
		year     INT    NOT NULL,

		PRIMARY KEY (id)
	)`,
}

type LocalDB struct {
	db *sql.DB
}

func NewLocalDB(config webapi.DBConfig) (*LocalDB, error) {
	var db, err = sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, err
	}

	return &LocalDB{db}, nil
}

func (local *LocalDB) Close() error {
	return local.db.Close()
}

func (local *LocalDB) Init() error {
	for i, ddl := range ddls {
		if _, err := local.db.Exec(ddl); err != nil {
			return fmt.Errorf("init %d: %w", i, err)
		}
	}

	return nil
}

func (local *LocalDB) Test() error {
	return local.withTempTx(context.Background(), func(db sqlDB) error {
		var (
			dbIDs = []int64{1, 2, 3, 4, 5}
			xcIDs = []int64{2, 4}
			xpIDs = []int64{1, 3, 5}
		)

		for i, id := range dbIDs {
			if _, err := db.Exec(`INSERT INTO audit.audit (id, num, month, year) VALUES ($1, $2, $3, $4)`, id, i+1, 3, 2022); err != nil {
				return fmt.Errorf("test INSERT %d: %w", i, err)
			}
		}

		var rows, err = db.Query(`SELECT id FROM audit.audit WHERE id != ALL($1)`, pq.Array(xcIDs))
		if err != nil {
			return fmt.Errorf("test SELECT: %w", err)
		}
		defer func() {
			_ = rows.Close()
		}()

		var gtIDs []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				return fmt.Errorf("test Scan: %w", err)
			}
			gtIDs = append(gtIDs, id)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		if len(xpIDs) != len(gtIDs) {
			return fmt.Errorf("len mismatch: expected %d, got %d", len(xcIDs), len(gtIDs))
		}

		for i := range xpIDs {
			if xpIDs[i] != gtIDs[i] {
				return fmt.Errorf("item mismatch at position %d: expceted %d, got %d", i, xpIDs[i], gtIDs[i])
			}
		}

		return nil
	}, nil)
}

func (local *LocalDB) withTempTx(ctx context.Context, fn func(sqlDB) error, opts *sql.TxOptions) (err error) {
	var tx *sql.Tx

	tx, err = local.db.BeginTx(ctx, opts)
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
			// all good, cleanup
			err = tx.Rollback()
		}
	}()

	err = fn(tx)
	return err
}

type (
	sqlDB interface {
		Exec(query string, params ...interface{}) (sql.Result, error)
		Query(query string, params ...interface{}) (*sql.Rows, error)
	}

	txSqlDB interface {
		WithTx(ctx context.Context, fn func(db sqlDB) error, opts *sql.TxOptions) error
	}
)
