package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/repository"
)

// По доступности соединения к БД проверяем готовность к работе
func handlerReadinessCheck(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	if rep != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	fmt.Fprintf(w, consts.TestReponseSuccess)
}

func handlerGetEntity(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	id, ok := parseEntityID(w, r)
	if !ok {
		return
	}
	result, err := rep.GetByID(id)
	if err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	tools.JSONRespOk(w, result)
}

func handlerUpdateEntity(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	id, ok := parseEntityID(w, r)
	if !ok {
		return
	}
	var item model.EntityCmd

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	defer r.Body.Close()

	err := rep.Update(id, &item)
	if err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	tools.JSONRespOk(w, consts.MakeSuccessResult())
}

func handlerDeleteEntity(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	id, ok := parseEntityID(w, r)
	if !ok {
		return
	}
	err := rep.Delete(id)
	if err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	tools.JSONRespOk(w, consts.MakeSuccessResult())
}

func handlerCreateEntity(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	var item model.EntityCmd

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	defer r.Body.Close()

	id, err := rep.Create(&item)
	if err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}

	result := consts.MakeSuccessResult()
	result[model.KeyNameID] = strconv.Itoa(int(id))
	tools.JSONRespOk(w, result)
}

func handlerGetEntities(rep repository.IEntityRepository, w http.ResponseWriter, r *http.Request) {
	result, err := rep.GetAll()
	if err != nil {
		makeRespInternalServerErr(w, r, err)
		return
	}
	tools.JSONRespOk(w, result)
}
