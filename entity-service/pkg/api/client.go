package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/model"

	"github.com/pkg/errors"
)

const (
	// ESvcMainPathURL - entity rest service relation path of main action
	ESvcMainPathURL = "/entities"

	logMarker = "ESvc: "
)

// EntityRestClient - construct executor of entity rest service
type EntityRestClient struct {
	TxtBaseURL string
	baseURL    *url.URL
}

// NewEntityRestClient - ctor of NewEntityRestClient
func NewEntityRestClient(txtBaseURL string) (*EntityRestClient, error) {
	u, err := url.Parse(txtBaseURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Entity Service Url is invalid [%s]", txtBaseURL)
	}
	return &EntityRestClient{
		TxtBaseURL: txtBaseURL,
		baseURL:    u,
	}, nil
}

func (x EntityRestClient) makeURL(action string) string {
	u := *x.baseURL
	u.Path = path.Join(u.Path, action)
	return u.String()
}
func (x EntityRestClient) makeURLWithID(action string, id int) string {
	u := *x.baseURL
	u.Path = path.Join(u.Path, action, strconv.Itoa(id))
	return u.String()
}

func checkSuccessResult(data map[string]string, respBody io.Reader) (bool, error) {
	v, present := data[consts.KeyNameResult]
	if !present {
		msg := fmt.Sprintf("Result field with key [%s] is not found. DataMap: %+v", consts.KeyNameResult, data)
		body, err2 := ioutil.ReadAll(respBody)
		if err2 == nil {
			msg += fmt.Sprintf(". RespBody: [%s]", string(body))
		}
		return false, errors.Errorf(msg)
	} else if v != consts.ValResultSuccess {
		return false, errors.Errorf("Failure  result. %+v", data)
	}
	return true, nil
}

func processError(msg string, err error) error {
	tools.Errorfln("%s. %s", msg, err.Error())
	return errors.Wrap(err, msg)
}

// MakeReadinessURL - make entity rest service url of readiness
func (x EntityRestClient) MakeReadinessURL() (path string, method string) {
	path = x.makeURL(consts.ReadinessPathURL)
	method = http.MethodGet
	return
}

// MakeGetAllEntitiesURL - make entity rest service url of GetAllEntities
func (x EntityRestClient) MakeGetAllEntitiesURL() (path string, method string) {
	path = x.makeURL(ESvcMainPathURL)
	method = http.MethodGet
	return
}

// MakeCreateURL - make entity rest service url of Create Entity
func (x EntityRestClient) MakeCreateURL() (path string, method string) {
	path = x.makeURL(ESvcMainPathURL)
	method = http.MethodPost
	return
}

// MakeGetEntityURL - make entity rest service url of Get Entity
func (x EntityRestClient) MakeGetEntityURL(id int) (path string, method string) {
	path = x.makeURLWithID(ESvcMainPathURL, id)
	method = http.MethodGet
	return
}

// MakeUpdateURL - make entity rest service url of Update Entity
func (x EntityRestClient) MakeUpdateURL(id int) (path string, method string) {
	path = x.makeURLWithID(ESvcMainPathURL, id)
	method = http.MethodPut
	return
}

// MakeDeleteURL - make entity rest service url of Delete Entity
func (x EntityRestClient) MakeDeleteURL(id int) (path string, method string) {
	path = x.makeURLWithID(ESvcMainPathURL, id)
	method = http.MethodDelete
	return
}

// ExecReadiness - get StatusCode  for Readiness action
func (x EntityRestClient) ExecReadiness() (int, error) {
	u, m := x.MakeReadinessURL()

	msg := fmt.Sprintf("%s(%s %s): %s", logMarker, m, u, "ExecReadiness")
	tools.Debugln(msg)

	resp, err := http.Get(u)
	if err != nil {
		err = processError(msg, err)
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// ExecGetAllEntities - get content for GetAllEntities action
func (x EntityRestClient) ExecGetAllEntities() ([]model.EntityQuery, error) {
	var data []model.EntityQuery

	u, m := x.MakeGetAllEntitiesURL()

	msg := fmt.Sprintf("%s(%s %s): %s", logMarker, m, u, "ExecGetAllEntities")
	tools.Debugln(msg)

	resp, err := http.Get(u)
	if err != nil {
		err = processError(msg, err)
		return data, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		err = processError(msg, err)
		return data, err
	}

	return data, nil
}

// ExecGetEntity - get content for GetEntity action
func (x EntityRestClient) ExecGetEntity(id int) (*model.EntityQuery, error) {
	var data *model.EntityQuery

	u, m := x.MakeGetEntityURL(id)

	msg := fmt.Sprintf("%s(%s %s): %s", logMarker, m, u, "ExecGetEntity")
	tools.Debugln(msg)

	resp, err := http.Get(u)
	if err != nil {
		err = processError(msg, err)
		return data, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		err = processError(msg, err)
		return data, err
	}

	return data, nil
}

// ExecDeleteEntity - get content for DeleteEntity action
func (x EntityRestClient) ExecDeleteEntity(id int) error {

	u, m := x.MakeDeleteURL(id)

	msg := fmt.Sprintf("%s(%s %s): %s", logMarker, m, u, "ExecDeleteEntity")
	tools.Debugln(msg)

	client := &http.Client{}
	req, err := http.NewRequest(m, u, nil)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	defer resp.Body.Close()

	var data map[string]string
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	_, err = checkSuccessResult(data, resp.Body)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	return nil
}

// ExecCreateEntity - get content for CreateEntity action
func (x EntityRestClient) ExecCreateEntity(cmd model.EntityCmd) (int, error) {

	u, m := x.MakeCreateURL()

	msg := fmt.Sprintf("%s(%s %s): %s. Cmd:{%+v}", logMarker, m, u, "ExecCreateEntity", cmd)
	tools.Debugln(msg)

	inData, err := json.Marshal(cmd)
	if err != nil {
		err = processError(msg, err)
		return int(model.InvalidEntityID), err
	}

	client := &http.Client{}
	req, err := http.NewRequest(m, u, bytes.NewBuffer(inData))
	if err != nil {
		err = processError(msg, err)
		return int(model.InvalidEntityID), err
	}
	resp, err := client.Do(req)
	if err != nil {
		err = processError(msg, err)
		return int(model.InvalidEntityID), err
	}
	defer resp.Body.Close()

	var data map[string]string
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		err = processError(msg, err)
		return int(model.InvalidEntityID), err
	}

	_, err = checkSuccessResult(data, resp.Body)
	if err != nil {
		err = processError(msg, err)
		return int(model.InvalidEntityID), err
	}

	idText, present := data[model.KeyNameID]
	if !present {
		return int(model.InvalidEntityID), fmt.Errorf("Can't find [%s] field in the result [%#v]",
			model.KeyNameID, data)
	}
	id, err := strconv.Atoi(idText)
	if err != nil {
		return int(model.InvalidEntityID), err
	}
	return id, nil
}

// ExecUpdateEntity - get content for UpdateEntity action
func (x EntityRestClient) ExecUpdateEntity(id int, cmd model.EntityCmd) error {

	u, m := x.MakeUpdateURL(id)

	msg := fmt.Sprintf("%s(%s %s): %s. Cmd:{%+v}", logMarker, m, u, "ExecUpdateEntity", cmd)
	tools.Debugln(msg)

	inData, err := json.Marshal(cmd)
	if err != nil {
		err = processError(msg, err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(m, u, bytes.NewBuffer(inData))
	if err != nil {
		err = processError(msg, err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	defer resp.Body.Close()

	var data map[string]string
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		err = processError(msg, err)
		return err
	}
	_, err = checkSuccessResult(data, resp.Body)
	if err != nil {
		err = processError(msg, err)
		return err
	}

	return nil
}
