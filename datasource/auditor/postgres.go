package auditor

import (
	"database/sql"

	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) datamodel.Auditors {
	return &postgres{
		db: db,
	}
}

func (pg *postgres) FindOne(id datamodel.AuditorID) (_ *datamodel.Auditor, err error) {
	const query = `
		SELECT
			first_name,
			last_name,
			position_id
		FROM
			audit.auditors
		WHERE
			id = $1
	`

	var a = datamodel.Auditor{ID: id}
	err = pg.db.QueryRow(query, id).Scan(
		&a.FirstName,
		&a.LastName,
		&a.PositionID)

	if err == sql.ErrNoRows {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (pg *postgres) FindAll() (list []datamodel.Auditor, err error) {
	const query = `
		SELECT
			id,
			first_name,
			last_name,
			position_id
		FROM
			audit.auditors
	`

	var rows *sql.Rows
	rows, err = pg.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var a datamodel.Auditor
		err = rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.PositionID)

		if err != nil {
			return nil, err
		}

		list = append(list, a)
	}

	return list, rows.Err()
}

func (pg *postgres) Create(auditor datamodel.Auditor) (id datamodel.AuditorID, err error) {
	const query = `
		INSERT INTO audit.auditors (
			first_name,
			last_name,
			position_id
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = pg.db.QueryRow(query,
		auditor.FirstName,
		auditor.LastName,
		auditor.PositionID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
