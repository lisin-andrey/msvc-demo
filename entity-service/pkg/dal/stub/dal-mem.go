package stub

import (
	"sync"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

	"github.com/pkg/errors"
)

// EntityStubDal - implementation of IEntityDal for STUB
type EntityStubDal struct {
	muID     sync.Mutex
	lastID   int32
	entities *mapEntity
}

// NewEntityStubDal - ctor of EntityStubDal
func NewEntityStubDal() (*EntityStubDal, error) {
	return &EntityStubDal{
		lastID:   0,
		entities: newMapEntity(),
	}, nil
}

func (dal *EntityStubDal) nextID() int32 {
	dal.muID.Lock()
	defer dal.muID.Unlock()
	dal.lastID++
	return dal.lastID
}

// Close - see IEntityDal
func (dal *EntityStubDal) Close() {
	if dal.entities != nil {
		dal.entities = nil
		dal.lastID = 0
	}
}

// Create - see IEntityDal
func (dal *EntityStubDal) Create(v model.EntityCmd) (int32, error) {

	//check Name unique. Implementation for test purpose only
	items := dal.entities.getAll()
	for _, it := range items {
		if it.Name == v.Name {
			return model.InvalidEntityID, errors.Errorf("Try insert record with duplicated Name [%s]", v.Name)
		}
	}

	id := dal.nextID()
	eq := model.EntityQuery{
		ID:           id,
		Name:         v.Name,
		Descr:        v.Descr,
		Created:      v.LastUpdated,
		LastUpdated:  v.LastUpdated,
		LastOperator: v.LastOperator,
	}
	dal.entities.add(eq)
	return id, nil
}

// Update - see IEntityDal
func (dal *EntityStubDal) Update(id int32, v model.EntityCmd) (bool, error) {
	eq, ok := dal.entities.get(id)
	if !ok {
		return false, errors.Errorf("Item with ID = %d is absent in the store", id)
	}
	eq.Descr = v.Descr
	eq.LastUpdated = v.LastUpdated
	eq.LastOperator = v.LastOperator

	err := dal.entities.set(eq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Delete - see IEntityDal
func (dal *EntityStubDal) Delete(id int32) (bool, error) {
	_, ok := dal.entities.get(id)
	if !ok {
		return false, nil
	}
	dal.entities.delete(id)
	return true, nil
}

// DeleteAll - see IEntityDal
func (dal *EntityStubDal) DeleteAll() error {
	dal.entities.deleteAll()
	return nil
}

// GetByID - see IEntityDal
func (dal *EntityStubDal) GetByID(id int32) (*model.EntityQuery, error) {
	eq, ok := dal.entities.get(id)
	if !ok {
		return nil, nil
	}
	return &eq, nil
}

// GetAll - see IEntityDal
func (dal *EntityStubDal) GetAll() ([]model.EntityQuery, error) {
	return dal.entities.getAll(), nil
}
