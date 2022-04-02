package datasource

import (
	"github.com/andreyskoskin/drvolodko2/datamodel"
	"github.com/andreyskoskin/drvolodko2/datasource/auditor"
	"github.com/andreyskoskin/drvolodko2/datasource/auditprogram"
)

type (
	InMemory struct {
		auditPrograms datamodel.AuditPrograms
		auditors      datamodel.Auditors
	}
)

func NewInMemory() *InMemory {
	return &InMemory{
		auditPrograms: auditprogram.NewInMemory(),
		auditors:      auditor.NewInMemory(),
	}
}

func (mem *InMemory) AuditPrograms() datamodel.AuditPrograms {
	return mem.auditPrograms
}

func (mem *InMemory) Auditors() datamodel.Auditors {
	return mem.auditors
}

func (mem *InMemory) Close() error {
	return nil
}
