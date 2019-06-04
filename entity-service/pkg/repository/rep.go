package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/config"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/dal"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/dal/stub"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/dal/postgres"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

	"github.com/pkg/errors"
)

const (
	// ProviderNamePostgres - the name of provider to create IEntityDal
	ProviderNamePostgres = "postgres"
	// ProviderNameMemory - the name of provider to create IEntityDal
	ProviderNameMemory = "memory"
)

// EntityRepository - entity repository implementation
type EntityRepository struct {
	dal dal.IEntityDal
}

// NewEntityRepository - Ctor of the EntityRepository
func NewEntityRepository(v dal.IEntityDal) (*EntityRepository, error) {
	if v == nil {
		return nil, errors.New("IEntityDal is undefined")
	}
	er := EntityRepository{dal: v}
	return &er, nil
}

// NewEntityRepositoryByConfig - Ctor of the EntityRepository
func NewEntityRepositoryByConfig(config *config.Config) (*EntityRepository, error) {
	if config.ProviderType == "" {
		return nil, errors.New("DB Provider is empty")
	}

	var dal dal.IEntityDal
	var err error

	switch strings.ToLower(config.ProviderType) {
	case ProviderNamePostgres:
		dal, err = postgres.NewEntityPgDalByConfig(config.PgdDal)
	case ProviderNameMemory:
		dal, err = stub.NewEntityStubDal()
	default:
		return nil, errors.Errorf("Unknown DB provider [%s]", config.ProviderType)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Can't create EntityRepository for provider %s", config.ProviderType)
	}

	return NewEntityRepository(dal)
}

func (r *EntityRepository) validate() {
	if r.dal == nil {
		panic("EntityRepository was not initialized correctly")
	}
}

// Close - desctructor of EntityRepository
func (r *EntityRepository) Close() {
	if r.dal != nil {
		r.dal.Close()
		r.dal = nil
	}
}

// Create - create entity
func (r *EntityRepository) Create(v *model.EntityCmd) (int32, error) {
	r.validate()

	var err error
	var id int32

	if err = v.Validate(); err != nil {
		return model.InvalidEntityID, fmt.Errorf("Creating entity is invalid. [%s]", v)
	}
	v.LastUpdated = time.Now()

	id, err = r.dal.Create(*v)
	if err != nil {
		return model.InvalidEntityID, errors.Wrap(err, fmt.Sprintf("Can't create entity [%s]", v.String()))
	}
	return id, nil
}

// Update  - update entity
func (r *EntityRepository) Update(id int32, v *model.EntityCmd) error {
	r.validate()

	var err error

	if err = v.ValidateForUpdate(); err != nil {
		return fmt.Errorf("Updating entity with id = %d is invalid. [%s]", id, v)
	}
	v.LastUpdated = time.Now()

	ok, err := r.dal.Update(id, *v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't update entity with id = %d [%s]", id, v.String()))
	}
	if !ok {
		return fmt.Errorf("Can't update entity with id = %d. Record is not found", id)
	}
	return nil
}

// Delete - delete entity
func (r *EntityRepository) Delete(id int32) error {
	r.validate()

	_, err := r.dal.Delete(id)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Can't delete entity with id = %d", id))
	}
	return nil
}

// GetByID - get entity by id
func (r *EntityRepository) GetByID(id int32) (*model.EntityQuery, error) {
	r.validate()

	v, err := r.dal.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Can't get entity with id = %d", id))
	}
	return v, nil
}

// GetAll - get all entites
func (r *EntityRepository) GetAll() ([]model.EntityQuery, error) {
	r.validate()

	lst, err := r.dal.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "Can't get all entities")
	}
	return lst, nil
}
