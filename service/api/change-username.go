package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
import (

	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

)
*/
func (rt *_router) changeUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parsing the parameters from the request
	userID := ps.ByName("userID")

	var username = structs.Username{}

	err := json.NewDecoder(r.Body).Decode(&username)

	if err != nil {
		// ^BadRequest va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while parsing the request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start a transaction
	// ^Internal Server Error va aggiunto all'openapi come possibile risposta
	err = rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the username is already taken
	result, err := rt.db.UserAvailable(username.Username)
	if result {
		// ^Forbidden va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Username already taken")
		httpErrorResponse(rt, w, "Username already taken", http.StatusForbidden)
		rt.db.Rollback()
		return
	}

	// If is free, change the username
	err = rt.db.ChangeUsername(userID, username.Username)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while changing username")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	// Commit the transaction
	err = rt.db.Commit()

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Username changed")

}
