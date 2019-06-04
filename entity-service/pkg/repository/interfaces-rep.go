package repository

import "github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

//IEntityRepository - repository to work with entity storage
type IEntityRepository interface {
	Close()
	Create(v *model.EntityCmd) (int32, error)
	Update(id int32, v *model.EntityCmd) error
	Delete(id int32) error
	GetByID(id int32) (*model.EntityQuery, error)
	GetAll() ([]model.EntityQuery, error)
}
