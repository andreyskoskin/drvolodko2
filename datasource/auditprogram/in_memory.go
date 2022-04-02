package auditprogram

import (
	"sync"

	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type inMemory struct {
	items sync.Map
}

func NewInMemory() datamodel.AuditPrograms {
	return &inMemory{}
}

func (mem *inMemory) FindOne(id datamodel.AuditProgramID) (*datamodel.AuditProgram, error) {
	var v, found = mem.items.Load(id)
	if !found {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	var program, ok = v.(datamodel.AuditProgram)
	if !ok {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	return &program, nil
}

func (mem *inMemory) Update(program datamodel.AuditProgram) (err error) {
	if _, found := mem.items.Load(program.ID); !found {
		return datamodel.ErrAuditProgramNotFound
	}

	mem.items.Store(program.ID, &program)
	return nil
}
