package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Parsing the json input into a string
	var userDetails structs.UserInfo

	err := json.NewDecoder(r.Body).Decode(&userDetails)
	_ = r.Body.Close()

	username := userDetails.User

	//Checking for errors during the deconding process
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Trying to login the user
	userID, err := rt.db.LoginUser(username)
	if err == sql.ErrNoRows {

		//If the user's not registered add it to the database
		newUserID, err := uuid.NewV4()
		userID = newUserID.String()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Operation: " + ctx.ReqUUID.String() + "can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = rt.db.RegisterUser(userID, username)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Operation: " + ctx.ReqUUID.String() + "can't add the user to the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	rt.baseLogger.Info("User " + username + " logged in with uuid " + userID)

	// Set the content type to JSON
	w.Header().Set("content-type", "application/json; charset=UTF-8")

	//Prepare the JSON
	err = json.NewEncoder(w).Encode(structs.AuthToken{Token: userID})

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
