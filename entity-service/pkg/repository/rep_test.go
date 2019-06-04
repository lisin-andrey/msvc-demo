package repository

import (
	"testing"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/dal"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/dal/stub"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepSuite struct {
	suite.Suite
	rep *EntityRepository
	dal dal.IEntityDal
}

func (r *RepSuite) SetupSuite() {

	r.dal, _ = stub.NewEntityStubDal()
	r.rep, _ = NewEntityRepository(r.dal)
}

func (r *RepSuite) SetupTest() {
	r.dal.DeleteAll()
}

func (r *RepSuite) TearDownSuite() {
	r.dal.Close()
}

func TestRepSuite(t *testing.T) {
	s := new(RepSuite)
	suite.Run(t, s)
}

func (r *RepSuite) TestGetAllEmpty() {
	v, err := r.rep.GetAll()
	if assert.NoError(r.T(), err) {
		assert.Equal(r.T(), 0, len(v), "Expected empty list")
	}
}

func checkCreate(r *RepSuite, cmd *model.EntityCmd) int32 {
	id, err := r.rep.Create(cmd)
	if assert.NoError(r.T(), err, "Create failed") {
		assert.True(r.T(), id > 0, "incorrect id. Must be > 0")
	}
	return id
}
func (r *RepSuite) TestCreate() {
	cmd := model.EntityCmd{
		Name:         "TestName",
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}
	id := checkCreate(r, &cmd)

	v, err := r.rep.GetByID(id)
	if assert.NoError(r.T(), err, "GetByID(%d) failed", id) {
		assert.Equal(r.T(), id, v.ID, "ID")
		assert.Equal(r.T(), cmd.Name, v.Name, "Name")
		assert.Equal(r.T(), cmd.Descr, v.Descr, "Descr")
		assert.Equal(r.T(), cmd.LastOperator, v.LastOperator, "LastOperator")
		assert.False(r.T(), v.Created.IsZero(), "Created")
		assert.Equal(r.T(), v.Created, v.LastUpdated, "Created != LastUpdated")
	}
}

func (r *RepSuite) TestInvalidDataForCreate() {
	cmd := model.EntityCmd{
		Name:         "",
		LastOperator: "User 1",
	}
	_, err := r.rep.Create(&cmd)
	assert.Error(r.T(), err, "Name must be filled")

	cmd = model.EntityCmd{
		Name:         "Test Name",
		LastOperator: "",
	}
	_, err = r.rep.Create(&cmd)
	assert.Error(r.T(), err, "LastOperator must be filled")

	cmd.LastOperator = "User 1"
	_, err = r.rep.Create(&cmd)
	assert.NoError(r.T(), err, "Now record must be inserted")
}

func (r *RepSuite) TestTryCreateDuplicate() {
	cmd := model.EntityCmd{
		Name:         "TestName",
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}
	id := checkCreate(r, &cmd)

	v, err := r.rep.GetByID(id)
	if assert.NoError(r.T(), err, "GetByID(%d) failed", id) {
		assert.Equal(r.T(), id, v.ID, "ID")
	}

	// Try add duplicate
	_, err = r.rep.Create(&cmd)
	assert.Error(r.T(), err, "2 records with the same Name mustn't be created")
}

func (r *RepSuite) TestUpdate() {
	nameOrig := "TestName"
	cmd := model.EntityCmd{
		Name:         nameOrig,
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}

	id := checkCreate(r, &cmd)

	cmd.Name = "Name mustn't be changed"
	cmd.Descr += "-changed"
	cmd.LastOperator += "User 2"
	err := r.rep.Update(id, &cmd)
	assert.NoError(r.T(), err, "Update(%d, ...) failed", id)

	v, err := r.rep.GetByID(id)
	if assert.NoError(r.T(), err, "GetByID(%d) failed", id) {
		assert.Equal(r.T(), id, v.ID, "ID")
		assert.Equal(r.T(), nameOrig, v.Name, "Name")
		assert.Equal(r.T(), cmd.Descr, v.Descr, "Descr")
		assert.Equal(r.T(), cmd.LastOperator, v.LastOperator, "LastOperator")
		assert.False(r.T(), v.Created.IsZero(), "Created")
		assert.True(r.T(), v.LastUpdated.After(v.Created), "Created < LastUpdated")
	}
}

func (r *RepSuite) TestDelete() {
	cmd := model.EntityCmd{
		Name:         "TestName",
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}

	id := checkCreate(r, &cmd)

	err := r.rep.Delete(id)
	assert.NoError(r.T(), err, "Delete(%d) failed", id)

	v, err := r.rep.GetByID(id)
	if assert.NoError(r.T(), err, "GetByID(%d) failed", id) {
		assert.True(r.T(), v == nil, "Entity with id %d was not deleted", id)
	}
}

func (r *RepSuite) TestGetAllNotEmpty() {
	cmds := []model.EntityCmd{
		model.EntityCmd{
			Name:         "TestName 1",
			Descr:        "Descr for 1",
			LastOperator: "User 1",
		},
		model.EntityCmd{
			Name:         "TestName 2",
			Descr:        "Descr for 1",
			LastOperator: "User 2",
		},
	}
	id1 := checkCreate(r, &cmds[0])
	id2 := checkCreate(r, &cmds[1])

	v1, err := r.rep.GetAll()
	if assert.NoError(r.T(), err, "GetAll failed") {
		assert.True(r.T(), len(v1) == 2, "Expected 2 items in the storage")
	}

	err = r.rep.Delete(id1)
	assert.NoError(r.T(), err, "Delete(%d) failed", id1)

	v2, err := r.rep.GetAll()
	if assert.NoError(r.T(), err, "GetAll failed") {
		assert.True(r.T(), len(v2) == 1, "Expected 1 items in the storage")
	}

	err = r.rep.Delete(id2)
	assert.NoError(r.T(), err, "Delete(%d) failed", id2)

	v3, err := r.rep.GetAll()
	if assert.NoError(r.T(), err, "GetAll failed") {
		assert.True(r.T(), len(v3) == 0, "Expected 0 items in the storage")
	}
}
