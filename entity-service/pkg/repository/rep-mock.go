package repository

import (
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

	"github.com/stretchr/testify/mock"
)

// MockEntityRepository - entity repository implementation
type MockEntityRepository struct {
	mock.Mock
}

// Close - see IEntityRepository
func (r *MockEntityRepository) Close() {
	r.Called()
}

// Create - see IEntityRepository
func (r *MockEntityRepository) Create(v *model.EntityCmd) (int32, error) {
	rets := r.Called(v)
	return rets.Get(0).(int32), rets.Error(1)
}

// Update  - see IEntityRepository
func (r *MockEntityRepository) Update(id int32, v *model.EntityCmd) error {
	rets := r.Called(id, v)
	return rets.Error(0)
}

// Delete - see IEntityRepository
func (r *MockEntityRepository) Delete(id int32) error {
	rets := r.Called(id)
	return rets.Error(0)
}

// GetByID - see IEntityRepository
func (r *MockEntityRepository) GetByID(id int32) (*model.EntityQuery, error) {
	rets := r.Called(id)
	return rets.Get(0).(*model.EntityQuery), rets.Error(1)
}

// GetAll - see IEntityRepository
func (r *MockEntityRepository) GetAll() ([]model.EntityQuery, error) {
	rets := r.Called()
	return rets.Get(0).([]model.EntityQuery), rets.Error(1)
}
