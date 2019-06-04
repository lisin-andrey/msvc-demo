package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
	"github.com/lisin-andrey/msvc-demo/common/pkg/web"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/repository"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

var (
	_mockRep     *repository.MockEntityRepository
	_fRegHandler func(func(repository.IEntityRepository, http.ResponseWriter, *http.Request)) http.Handler
)

type mockRepInstanceHandler struct{}

func (h *mockRepInstanceHandler) GetRepInstance() repository.IEntityRepository {
	return _mockRep
}

func initMock(t *testing.T) {
	_mockRep = new(repository.MockEntityRepository)
	_fRegHandler = func(f func(repository.IEntityRepository, http.ResponseWriter, *http.Request)) http.Handler {
		return &RestHandler{
			RepHandler:     &mockRepInstanceHandler{},
			RestHandleFunc: f,
		}
	}
}

func checkSuccessResult(t *testing.T, recorder *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, recorder.Code)

	var actual map[string]string
	err := json.NewDecoder(recorder.Body).Decode(&actual)
	if assert.NoError(t, err, "Can't decode body [%s]", recorder.Body.String()) {
		v, present := actual[consts.KeyNameResult]
		assert.True(t, (present || v == consts.ValResultSuccess), "Handler returned unexpected body", recorder.Body.String())
	}
}

func TestHandlerLivenessCheck(t *testing.T) {
	initMock(t)
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(web.HandlerLivenessCheck)
	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, consts.TestReponseSuccess, recorder.Body.String(), "Handler returned unexpected body")
}

func TestHandlerReadinessCheck(t *testing.T) {
	initMock(t)

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerReadinessCheck))
	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, consts.TestReponseSuccess, recorder.Body.String(), "Handler returned unexpected body")
}

func TestHandlerGetEntity(t *testing.T) {
	initMock(t)

	etalon := model.EntityQuery{
		ID:           1,
		Name:         "TestName",
		Descr:        "Descr for 1",
		Created:      time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		LastUpdated:  time.Date(2019, 1, 2, 3, 4, 6, 0, time.UTC),
		LastOperator: "User 1",
	}

	_mockRep.On("GetByID", int32(1)).Return(&etalon, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerGetEntity))
	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var actual model.EntityQuery
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if assert.NoError(t, err, "Can't decode body [%s]", recorder.Body.String()) {
		assert.Equal(t, etalon, actual, "Handler returned unexpected body")
	}

	_mockRep.AssertExpectations(t)
}

func TestHandlerGetEntities(t *testing.T) {
	initMock(t)

	etalon := model.EntityQuery{
		ID:           1,
		Name:         "TestName",
		Descr:        "Descr for 1",
		Created:      time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		LastUpdated:  time.Date(2019, 1, 2, 3, 4, 6, 0, time.UTC),
		LastOperator: "User 1",
	}

	_mockRep.On("GetAll").Return([]model.EntityQuery{etalon}, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerGetEntities))
	hf.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var res []model.EntityQuery
	err = json.NewDecoder(recorder.Body).Decode(&res)
	if assert.NoError(t, err, "Can't decode body [%s]", recorder.Body.String()) {
		assert.Equal(t, 1, len(res), "Handler returned unexpected count of entities")
		assert.Equal(t, etalon, res[0], "Handler returned unexpected body")
	}

	_mockRep.AssertExpectations(t)
}

func TestHandlerDeleteEntity(t *testing.T) {
	initMock(t)

	_mockRep.On("Delete", int32(1)).Return(nil).Once()

	req, err := http.NewRequest(http.MethodDelete, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerDeleteEntity))
	hf.ServeHTTP(recorder, req)

	checkSuccessResult(t, recorder)

	_mockRep.AssertExpectations(t)
}

func TestHandlerCreateEntity(t *testing.T) {
	initMock(t)

	cmd := model.EntityCmd{
		Name:         "TestName",
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		t.Errorf("Can't encode EntityCmd. %s", err.Error())
	}

	_mockRep.On("Create", &cmd).Return(int32(1), nil).Once()

	req, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(b)))
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerCreateEntity))
	hf.ServeHTTP(recorder, req)

	checkSuccessResult(t, recorder)

	_mockRep.AssertExpectations(t)
}

func TestHandlerUpdateEntity(t *testing.T) {
	initMock(t)

	cmd := model.EntityCmd{
		Name:         "TestName",
		Descr:        "Descr for 1",
		LastOperator: "User 1",
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		t.Errorf("Can't encode EntityCmd. %s", err.Error())
	}

	_mockRep.On("Update", int32(1), &cmd).Return(nil).Once()

	req, err := http.NewRequest(http.MethodPut, "", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(b)))
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	recorder := httptest.NewRecorder()
	hf := http.Handler(_fRegHandler(handlerUpdateEntity))
	hf.ServeHTTP(recorder, req)

	checkSuccessResult(t, recorder)

	_mockRep.AssertExpectations(t)
}
