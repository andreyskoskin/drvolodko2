package auditprogram

import (
	"database/sql"

	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) datamodel.AuditPrograms {
	return &postgres{
		db: db,
	}
}

func (pg *postgres) FindOne(id datamodel.AuditProgramID) (_ *datamodel.AuditProgram, err error) {
	const query = `
		SELECT
			number,
			month,
			status,
			auditor_id,
			audit_type_id,
			speciality_id
		FROM
			audit.programs
		WHERE
			id = $1
	`

	var ap = datamodel.AuditProgram{ID: id}
	err = pg.db.QueryRow(query, id).Scan(
		&ap.Number,
		&ap.Month,
		&ap.Status,
		&ap.AuditorID,
		&ap.AuditTypeID,
		&ap.SpecialityID,
	)

	if err == sql.ErrNoRows {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	if err != nil {
		return nil, err
	}

	return &ap, nil
}

func (pg *postgres) Update(program datamodel.AuditProgram) (err error) {
	const query = `
		UPDATE audit.programs SET
			number        = $2,
			month         = $3,
			status        = $4,
			auditor_id    = $5,
			audit_type_id = $6,
			speciality_id = $7
		WHERE
			id = $1
	`

	var result sql.Result
	result, err = pg.db.Exec(query,
		program.ID,
		program.Number,
		program.Month,
		program.Status,
		program.AuditorID,
		program.AuditTypeID,
		program.SpecialityID,
	)

	if err != nil {
		return err
	}

	var affected int64
	if affected, err = result.RowsAffected(); err != nil {
		return err
	}

	if affected == 0 {
		return datamodel.ErrAuditProgramNotFound
	}

	return nil
}
