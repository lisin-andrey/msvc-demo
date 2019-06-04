package cmd

import (
	"fmt"
	"net/http"

	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"
)

// func handlerTest(wd webData, w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "test")
// }

// По доступности соединения к Entity REST сервису проверяем готовность к работе
func handlerReadinessCheck(wd webData, w http.ResponseWriter, r *http.Request) {
	status, err := wd.entityRestClient.ExecReadiness()
	w.WriteHeader(status)

	msg := fmt.Sprintf("Entity Rest Service returned status %d", status)
	if err != nil {
		msg += fmt.Sprintf(". Error: %s", err.Error())
	}

	data := struct{ Message string }{
		Message: msg,
	}
	preventCache(w)
	err = wd.tmpltTest.Execute(w, data)
	checkErrResp(wd, err, "Can't execute tmpltTest", w)
}

func _handlerGetEntities(wd webData, w http.ResponseWriter, r *http.Request, msgs []string) {
	data, err := wd.entityRestClient.ExecGetAllEntities()

	if checkErrRespFromEntitySvc(wd, err, "GetAllEntities", w, r, false) {
		return
	}

	viewData := struct {
		Entities []model.EntityQuery
		Messages []string
	}{
		Entities: data,
		Messages: msgs,
	}

	preventCache(w)
	err = wd.tmpltEntities.Execute(w, viewData)
	checkErrResp(wd, err, "Can't execute tmpltEntities", w)
}

func handlerGetAllEntites(wd webData, w http.ResponseWriter, r *http.Request) {
	msgs := parseMessagesFromCookies(w, r)
	_handlerGetEntities(wd, w, r, msgs)
}

func handlerGetEntity(wd webData, w http.ResponseWriter, r *http.Request) {
	id, ok := parsePrmID(wd, w, r)
	if !ok {
		return
	}

	data, err := wd.entityRestClient.ExecGetEntity(id)
	if checkErrRespFromEntitySvc(wd, err, fmt.Sprintf("GetEntity(%d)", id), w, r, true) {
		return
	}

	if data != nil {
		preventCache(w)
		err = wd.tmpltEntity.Execute(w, *data)
		checkErrResp(wd, err, "Can't execute tmpltEntity", w)
	} else {
		redirectWithMessages(w, r, pathGetAllEntites, http.StatusFound, []string{fmt.Sprintf("Entity with ID=%d is not found", id)})
		// _handlerGetEntities(wd, w, r, []string{fmt.Sprintf("Entity with ID=%d is not found", id)})
	}
}

func handlerDelete(wd webData, w http.ResponseWriter, r *http.Request) {
	id, ok := parsePrmID(wd, w, r)
	if !ok {
		return
	}

	err := wd.entityRestClient.ExecDeleteEntity(id)
	if checkErrRespFromEntitySvc(wd, err, fmt.Sprintf("DeleteEntity(%d)", id), w, r, true) {
		return
	}

	http.Redirect(w, r, pathGetAllEntites, http.StatusFound)
}

func handlerCreateDisplay(wd webData, w http.ResponseWriter, r *http.Request) {
	err := wd.tmpltCreate.Execute(w, struct{}{})
	checkErrResp(wd, err, "Can't execute tmpltCreate", w)
}

func handlerCreateExec(wd webData, w http.ResponseWriter, r *http.Request) {
	pCmd := parseEntityCmd(wd, w, r)
	if pCmd == nil {
		return
	}

	_, err := wd.entityRestClient.ExecCreateEntity(*pCmd)
	if checkErrRespFromEntitySvc(wd, err, "CreateEntity", w, r, true) {
		return
	}

	http.Redirect(w, r, pathGetAllEntites, http.StatusFound)
}

func handlerUpdateDisplay(wd webData, w http.ResponseWriter, r *http.Request) {
	id, ok := parsePrmID(wd, w, r)
	if !ok {
		return
	}

	data, err := wd.entityRestClient.ExecGetEntity(id)
	if checkErrRespFromEntitySvc(wd, err, fmt.Sprintf("GetEntity(%d)", id), w, r, true) {
		return
	}

	if data != nil {
		preventCache(w)
		err = wd.tmpltUpdate.Execute(w, *data)
		checkErrResp(wd, err, "Can't execute tmpltUpdate", w)
	} else {
		redirectWithMessages(w, r, pathGetAllEntites, http.StatusFound, []string{fmt.Sprintf("Entity with ID=%d is not found. Operation was canceled", id)})
		//_handlerGetEntities(wd, w, r, []string{fmt.Sprintf("Entity with ID=%d is not found. Operation was canceled", id)})
	}
}

func handlerUpdateExec(wd webData, w http.ResponseWriter, r *http.Request) {
	id, ok := parsePrmID(wd, w, r)
	if !ok {
		return
	}
	pCmd := parseEntityCmd(wd, w, r)
	if pCmd == nil {
		return
	}

	err := wd.entityRestClient.ExecUpdateEntity(id, *pCmd)
	if checkErrRespFromEntitySvc(wd, err, fmt.Sprintf("UpdateEntity(%d)", id), w, r, true) {
		return
	}

	http.Redirect(w, r, pathGetAllEntites, http.StatusFound)
}
