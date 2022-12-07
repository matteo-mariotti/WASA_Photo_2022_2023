package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// changeUsername parses the request extracting the user id and the new username then, after checking if the new username is available, it changes the username
func (rt *_router) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userID := ps.ByName("userID")

	var username = structs.Username{}

	err := json.NewDecoder(r.Body).Decode(&username)

	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while parsing the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start a transaction
	err = rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the username is already taken
	result, err := rt.db.UserAvailable(username.Username)
	if result {
		rt.baseLogger.WithError(err).Error("Username already taken")
		httpErrorResponse(rt, w, "Username already taken", http.StatusConflict)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// If is free, change the username
	err = rt.db.ChangeUsername(userID, username.Username)

	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while changing username")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Commit the transaction
	err = rt.db.Commit()

	if err != nil {
		rt.baseLogger.WithError(err).Error("Unable to commit")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Username changed")

}
