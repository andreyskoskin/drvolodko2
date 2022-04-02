package datamodel

import (
	"errors"
)

var ErrAuditorNotFound = errors.New("auditor not found")

type (
	Auditors interface {
		FindOne(id AuditorID) (*Auditor, error)
		FindAll() ([]Auditor, error)
		Create(auditor Auditor) (AuditorID, error)
	}

	AuditorID int64

	Auditor struct {
		ID         AuditorID `json:"id,omitempty"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		PositionID int64     `json:"position_id"`
	}
)
