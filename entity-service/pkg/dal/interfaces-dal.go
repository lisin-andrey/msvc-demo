package dal

import (
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
)

//IEntityDal - CRUD interface to work with entity storage
type IEntityDal interface {
	Create(item model.EntityCmd) (int32, error)
	Update(id int32, item model.EntityCmd) (bool, error)
	Delete(id int32) (bool, error)
	DeleteAll() error
	GetByID(id int32) (*model.EntityQuery, error)
	GetAll() ([]model.EntityQuery, error)

	Close()
}
