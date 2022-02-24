package audit

type (
	AuditorsGroup struct {
		DivisionID int64 `json:"division_id"`
		AuditorID  int64 `json:"auditor_id"`
		SpecID     int64 `json:"spec_id"`
		CategoryID int64 `json:"category_id"`
	}

	AuditorsGroups []AuditorsGroup
)

func (ags AuditorsGroups) len() int {
	return len(ags)
}

func (ags AuditorsGroups) insertInto(db SqlDB, auditID int64, i int) error {
	return ags[i].insertInto(db, auditID)
}

func (ags AuditorsGroups) deleteAllFrom(db SqlDB, auditID int64) (err error) {
	_, err = db.Exec(`DELETE FROM audit.divisions WHERE audit_id = $1`, auditID)
	return err
}

func (ag AuditorsGroup) insertInto(db SqlDB, auditID int64) (err error) {
	_, err = db.Exec(`
		INSERT INTO audit.auditors_group (audit_id, auditor_id, division_id, spec_id, category_id)
		VALUES ($1, $2, $3, $4, $5)
	`,
		auditID, ag.AuditorID, ag.DivisionID, ag.SpecID, ag.CategoryID,
	)
	return err
}
