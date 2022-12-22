package api

import (
	"WASA_Photo/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var userList []string

	// Get username from query string
	username := r.URL.Query().Get("username")

	// Get page number from query string
	page := r.URL.Query().Get("page")

	// Convert page number to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		// If the page number is not a number, set it to 1
		pageInt = 0
	}

	if pageInt < 0 {
		pageInt = 0
	}

	// Get the list of users
	userList, err = rt.db.GetUsers(username, pageInt*30, ctx.Token)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Error("No more users are available with this prefix: " + username)
		httpErrorResponse(rt, w, "404 Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while getting users")
		httpErrorResponse(rt, w, "Internal Sever Error", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(userList)

	if err != nil {
		rt.baseLogger.Error("Error while encoding the response")
		httpErrorResponse(rt, w, "Internal Sever Error", http.StatusBadRequest)
		return
	}

	rt.baseLogger.Info("Get users request completed")

}
