package datamodel

import (
	"errors"
	"time"
)

var ErrAuditProgramNotFound = errors.New("audit program not found")

type (
	AuditPrograms interface {
		FindOne(id AuditProgramID) (*AuditProgram, error)
		Update(program AuditProgram) error
	}

	AuditProgramID int

	AuditProgramStatus string

	AuditProgram struct {
		ID           AuditProgramID     `json:"id,omitempty"`
		Number       int                `json:"program_number"`
		Month        time.Month         `json:"month"`
		Status       AuditProgramStatus `json:"status"`
		AuditorID    AuditorID          `json:"auditor_id"`
		AuditTypeID  int64              `json:"audit_type_id"`
		SpecialityID int64              `json:"speciality_id"`
	}
)

const (
	AuditProgramStatusDraft      AuditProgramStatus = "draft"
	AuditProgramStatusInProgress AuditProgramStatus = "in_progress"
	AuditProgramStatusFinished   AuditProgramStatus = "finished"
)
