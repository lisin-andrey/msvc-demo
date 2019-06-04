package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
	"github.com/lisin-andrey/msvc-demo/common/pkg/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	pathCreate        = "/create"
	pathUpdate        = "/update"
	pathDelete        = "/delete"
	pathGetAllEntites = "/entities"
	pathGetEntity     = "/entity"

	queryPrmNameID         = "id"
	queryPrmNameTempCookie = "t"
)

// NewRouter - ctor for router
func NewRouter(h func(WebHandleFunc) http.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(consts.HealthyPathURL, web.HandlerLivenessCheck).Methods(http.MethodGet)
	r.Handle(consts.ReadinessPathURL, h(handlerReadinessCheck)).Methods(http.MethodGet)
	r.Handle(consts.MetricsPathURL, promhttp.Handler())

	//r.Handle("/test", h(handlerTest)).Methods(http.MethodGet)

	r.Handle(pathCreate, h(handlerCreateDisplay)).Methods(http.MethodGet)
	r.Handle(pathCreate, h(handlerCreateExec)).Methods(http.MethodPost)
	r.Handle(pathUpdate, h(handlerUpdateDisplay)).Methods(http.MethodGet)
	r.Handle(pathUpdate, h(handlerUpdateExec)).Methods(http.MethodPost)
	r.Handle(pathGetAllEntites, h(handlerGetAllEntites)).Methods(http.MethodGet)
	r.Handle(pathGetEntity, h(handlerGetEntity)).Methods(http.MethodGet)
	r.Handle(pathDelete, h(handlerDelete)).Methods(http.MethodGet)
	return r
}
