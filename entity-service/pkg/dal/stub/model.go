package stub

import (
	"sync"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

	"github.com/pkg/errors"
)

type mapEntity struct {
	mu       sync.RWMutex
	internal map[int32]model.EntityQuery
}

func newMapEntity() *mapEntity {
	return &mapEntity{internal: make(map[int32]model.EntityQuery)}
}

func (m *mapEntity) get(id int32) (model.EntityQuery, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.internal[id]
	return v, ok
}

func (m *mapEntity) getAll() []model.EntityQuery {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.internal) == 0 {
		return []model.EntityQuery{}
	}
	result := make([]model.EntityQuery, 0, len(m.internal))
	for _, v := range m.internal {
		result = append(result, v)
	}
	return result
}

func (m *mapEntity) add(v model.EntityQuery) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, present := m.internal[v.ID]; present {
		return errors.Errorf("Item with ID = %d is present on the store", v.ID)
	}
	m.internal[v.ID] = v
	return nil
}

func (m *mapEntity) set(v model.EntityQuery) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, present := m.internal[v.ID]; present {
		m.internal[v.ID] = v
	} else {
		return errors.Errorf("Item with ID = %d is absent in the store", v.ID)
	}
	return nil
}

func (m *mapEntity) delete(id int32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.internal, id)
}

func (m *mapEntity) deleteAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k := range m.internal {
		delete(m.internal, k)
	}
}
