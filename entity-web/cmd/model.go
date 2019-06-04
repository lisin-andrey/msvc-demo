package cmd

import (
	"html/template"
	"net/http"
	"path"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"
	"github.com/lisin-andrey/msvc-demo/entity-service/pkg/api"

	"github.com/pkg/errors"
)

// webData - app data which is available for each work http handler
type webData struct {
	entityRestClient *api.EntityRestClient

	tmpltErr      *template.Template
	tmpltTest     *template.Template
	tmpltEntities *template.Template
	tmpltEntity   *template.Template
	tmpltCreate   *template.Template
	tmpltUpdate   *template.Template
}

// newWebData - ctor of WebData
func newWebData(esvcURL string) (*webData, error) {
	exec, err := api.NewEntityRestClient(esvcURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't create EntityRestExecutor(%s)", esvcURL)
	}
	w := &webData{entityRestClient: exec}
	err = w.initTemplates()

	if err != nil {
		return nil, err
	}

	return w, nil
}

func (x *webData) initTemplates() error {
	const pathTmpltDir = "./templates/"
	return x.initTemplatesFromDir(pathTmpltDir)
}

func (x *webData) initTemplatesFromDir(pathTmpltDir string) error {
	var err error

	f := func(file string) (*template.Template, error) {
		p := path.Join(pathTmpltDir, file)
		tools.Debugln("template path: ", p)
		t, e := template.ParseFiles(p)
		if e != nil {
			return nil, errors.Wrapf(e, "Failed to parse template [%s]", p)
		}
		return t, nil
	}

	procLst := []struct {
		t **template.Template
		f string
	}{
		{t: &x.tmpltErr, f: "err.html"},
		{t: &x.tmpltTest, f: "test.html"},
		{t: &x.tmpltEntities, f: "entities.html"},
		{t: &x.tmpltEntity, f: "entity.html"},
		{t: &x.tmpltCreate, f: "create.html"},
		{t: &x.tmpltUpdate, f: "update.html"},
	}
	for _, v := range procLst {
		*v.t, err = f(v.f)
		if err != nil {
			return err
		}
	}
	return nil
}

// WebHandleFunc - Used format of http HandleFuncs in the project
type WebHandleFunc func(webData, http.ResponseWriter, *http.Request)

// WebHandler - wrapper of http.Handler
type WebHandler struct {
	webData
	WebHandleFunc
}

// implement http.Handler interface
func (h *WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.WebHandleFunc(h.webData, w, r)
}
