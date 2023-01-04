package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// login is the handler for POST /session which allows a user to login or signup into the web application.
// It returns a session token that can be used to authenticate the user in subsequent requests
// If there's and error during the deconding of the request body, it replies with a 400 Bad Request, if there is an error while adding the user to the database or while encoding the response, it replies with a 500 Internal Server Error.
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the json input into a string
	var userDetails structs.UserInfo

	err := json.NewDecoder(r.Body).Decode(&userDetails)

	username := userDetails.User

	// Checking for errors during the deconding process
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if username == "" {
		rt.baseLogger.Warning("No username received")
		httpErrorResponse(rt, w, "No username received", http.StatusBadRequest)
		return
	}

	// Trying to login the user
	userID, err := rt.db.LoginUser(username)
	if errors.Is(err, sql.ErrNoRows) {

		// If the user's not registered add it to the database
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

	// Prepare the JSON
	err = json.NewEncoder(w).Encode(structs.AuthToken{Token: userID})

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Error enconding", http.StatusInternalServerError)
		return
	}
}
