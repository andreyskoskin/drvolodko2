package audit

import (
	"time"
)

type (
	Division struct {
		DivisionID   int64     `json:"division_id"`
		StartDate    time.Time `json:"start_date"`
		EndDate      time.Time `json:"end_date"`
		PlaceID      int64     `json:"place_id"`
		ResponseID   int64     `json:"response_id"`
		ProfessionID int64     `json:"profession_id"`
	}

	Divisions []Division
)

func (ds Divisions) len() int {
	return len(ds)
}

func (ds Divisions) insertInto(db SqlDB, auditID int64, i int) error {
	return ds[i].insertInto(db, auditID)
}

func (ds Divisions) deleteAllFrom(db SqlDB, auditID int64) (err error) {
	_, err = db.Exec(`DELETE FROM audit.divisions WHERE audit_id = $1`, auditID)
	return err
}

func (d Division) insertInto(db SqlDB, auditID int64) (err error) {
	_, err = db.Exec(`
		INSERT INTO audit.division (audit_id, division_id, start_date, end_date, place_id, response_id, profession_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		auditID, d.DivisionID, d.StartDate, d.EndDate, d.PlaceID, d.ResponseID, d.ProfessionID,
	)
	return err
}
