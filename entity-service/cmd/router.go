package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
	"github.com/lisin-andrey/msvc-demo/common/pkg/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	pathEntities = "/entities"
	suffixPathID = "/{id:[0-9]+}"
	urlPrmNameID = "id"
)

// NewRouter - ctor for router
func NewRouter(h func(RestHandleFunc) http.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(consts.HealthyPathURL, web.HandlerLivenessCheck).Methods(http.MethodGet)
	r.Handle(consts.ReadinessPathURL, h(handlerReadinessCheck)).Methods(http.MethodGet)
	r.Handle(consts.MetricsPathURL, promhttp.Handler())

	r.Handle(pathEntities+suffixPathID, h(handlerGetEntity)).Methods(http.MethodGet)
	r.Handle(pathEntities+suffixPathID, h(handlerUpdateEntity)).Methods(http.MethodPut)
	r.Handle(pathEntities+suffixPathID, h(handlerDeleteEntity)).Methods(http.MethodDelete)
	r.Handle(pathEntities, h(handlerCreateEntity)).Methods(http.MethodPost)

	r.Handle(pathEntities, h(handlerGetEntities)).Methods(http.MethodGet)

	// r.Use(loggingMiddleware)
	return r
}

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tools.Debugln(r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	})
// }
