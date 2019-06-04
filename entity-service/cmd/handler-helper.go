package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lisin-andrey/msvc-demo/common/pkg/tools"

	"github.com/gorilla/mux"
)

const (
	msgInvalidEntityID = "Invalid EntityId."
)

func makeRespBadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	tools.JSONRespError(w, http.StatusBadRequest,
		msg+" Url: "+tools.GetURL(r))
}

func makeRespInternalServerErr(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error. Url: [%s]. %s", tools.GetURL(r), err.Error())
	tools.JSONRespError(w, http.StatusInternalServerError,
		fmt.Sprintf("%s", err.Error()))
}

func parseEntityID(w http.ResponseWriter, r *http.Request) (int32, bool) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars[urlPrmNameID], 10, 64)
	if err != nil {
		makeRespBadRequest(w, r, msgInvalidEntityID)
		return 0, false
	}
	return int32(id), true
}
