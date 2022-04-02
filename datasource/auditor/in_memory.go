package auditor

import (
	"sync"
	"sync/atomic"

	"github.com/andreyskoskin/drvolodko2/datamodel"
)

type inMemory struct {
	ids   *int64
	items sync.Map
}

func NewInMemory() datamodel.Auditors {
	return &inMemory{
		ids: new(int64),
	}
}

func (mem *inMemory) FindOne(id datamodel.AuditorID) (*datamodel.Auditor, error) {
	var v, found = mem.items.Load(id)
	if !found {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	var auditor, ok = v.(datamodel.Auditor)
	if !ok {
		return nil, datamodel.ErrAuditProgramNotFound
	}

	return &auditor, nil
}

func (mem *inMemory) FindAll() ([]datamodel.Auditor, error) {
	var list []datamodel.Auditor

	mem.items.Range(func(key, value any) bool {
		if auditor, ok := value.(datamodel.Auditor); ok {
			list = append(list, auditor)
		}
		return true
	})

	return list, nil
}

func (mem *inMemory) Create(auditor datamodel.Auditor) (datamodel.AuditorID, error) {
	auditor.ID = datamodel.AuditorID(atomic.AddInt64(mem.ids, 1))
	mem.items.Store(auditor.ID, auditor)
	return auditor.ID, nil
}
