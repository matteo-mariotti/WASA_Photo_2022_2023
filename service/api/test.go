package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	s := r.Header.Get("Authorization")
	rt.baseLogger.Info(s)
	splitToken := strings.Split(s, "Bearer ")
	s = splitToken[1]
	rt.baseLogger.Info(s)

}
