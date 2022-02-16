package audit

import (
	"context"
	"database/sql"
)

// external dependencies
type (
	SqlDB interface {
		Exec(query string, params ...interface{}) (sql.Result, error)
	}

	TxSqlDB interface {
		WithTx(ctx context.Context, fn func(db SqlDB) error, opts *sql.TxOptions) error
	}
)

type Audit struct {
	ID          int64  `json:"id,omitempty"`
	Num         string `json:"num"`
	Month       string `json:"month"`
	Year        int    `json:"year"`
	ThorAuditID int64  `json:"thor_audit_id"`
	AuditorID   int64  `json:"auditor_id"`
	IsSpecial   bool   `json:"is_special"`

	AuditorsGroups AuditorsGroups `json:"auditors_groups"`
	Divisions      Divisions      `json:"divisions"`
}

func (audit *Audit) SaveTo(db TxSqlDB, ctx context.Context) error {
	return db.WithTx(ctx, func(db SqlDB) (err error) {

		if err = audit.saveItems(db, audit.AuditorsGroups); err != nil {
			return err
		}

		if err = audit.saveItems(db, audit.Divisions); err != nil {
			return err
		}

		return audit.updateIn(db)
	}, nil)
}

type items interface {
	len() int
	insertInto(db SqlDB, parentID int64, i int) error
	deleteAllFrom(db SqlDB, parentID int64) error
}

func (audit *Audit) saveItems(db SqlDB, c items) error {
	if err := c.deleteAllFrom(db, audit.ID); err != nil {
		return err
	}

	var n = c.len()
	for i := 0; i < n; i++ {
		if err := c.insertInto(db, audit.ID, i); err != nil {
			return err
		}
	}
	return nil
}

func (audit *Audit) updateIn(db SqlDB) (err error) {
	_, err = db.Exec(`
		UPDATE audit.audit SET
			year  = $2,
			month = $3,
			num   = $4,
			thor_audit_id = $5,
			auditor_id    = $6,
			is_special    = $7
		WHERE id = $1
	`,
		audit.ID, audit.Year, audit.Month, audit.Num, audit.ThorAuditID, audit.AuditorID, audit.IsSpecial,
	)
	return err
}
