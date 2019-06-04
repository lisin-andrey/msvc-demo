package web

import (
	"fmt"
	"net/http"

	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
)

// HandlerLivenessCheck - Handler for livenessProbe
func HandlerLivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, consts.TestReponseSuccess)
}
