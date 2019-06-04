package tools

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lisin-andrey/msvc-demo/common/pkg/consts"
)

// GetURL - extract url from http.Request
func GetURL(r *http.Request) string {
	return fmt.Sprintf("%s%s", r.Host, r.URL.Path)
}

// SetJSONRespHeaders - set headers for JSON response
func SetJSONRespHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set(consts.HTTPHdrNameContentType, consts.HTTPContentTypeJSON)
	w.WriteHeader(statusCode)
}

// JSONRespOk - make JSON response with code 200
func JSONRespOk(w http.ResponseWriter, body interface{}) {
	SetJSONRespHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

// JSONRespError - make JSON response for error message
func JSONRespError(w http.ResponseWriter, code int, message string) {
	SetJSONRespHeaders(w, code)

	body := consts.MakeFailureResult(message)
	json.NewEncoder(w).Encode(body)
}
