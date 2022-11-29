package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var profileResponse structs.Profile

	// Parsing the parameters from the request
	userProfile := ps.ByName("userID")

	// Check if the user I'm trying to look at has blocked me
	isBanned, err := rt.db.IsBanned(userProfile, ctx.Token)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userProfile + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to get the profile: userB has banned userA. userB: " + userProfile + " userA: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot get the profile of a user that has blocked you", http.StatusForbidden)
		return
	}

	// If not, get the profile
	username, err := rt.db.GetName(userProfile)

	if err == sql.ErrNoRows {
		//^Aggiungere NotFound come possibile risposta all'openapi
		rt.baseLogger.Error("This user does not exist!" + userProfile)
		httpErrorResponse(rt, w, "UserID is not valid", http.StatusNotFound)
		return
	} else if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the username of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the username
	profileResponse.Username = username

	// Get the number of followers
	follower, err := rt.db.GetFollowerNumber(userProfile)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the number of followers of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of followers
	profileResponse.Follower = follower

	// Get the number of following
	following, err := rt.db.GetFollowingNumber(userProfile)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the number of following of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of following
	profileResponse.Following = following

	// Get the number of photos
	photoNumber, err := rt.db.GetPhotosNumber(userProfile)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the number of photos of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of photos
	profileResponse.PhotoNumber = photoNumber
	/*
		// Get the photos (using the offset)
		photos, err := rt.db.GetPhotos(userProfile, 0)

		if err != nil {
			//^Aggiungere InternalServerError come possibile risposta all'openapi
			rt.baseLogger.Error("Error while getting the photos of user " + userProfile)
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the photos
		profileResponse.Photo = photos
	*/
	//Prepare the JSON
	err = json.NewEncoder(w).Encode(profileResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	// If everything went well, return 204
	w.WriteHeader(http.StatusOK)

	//Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the profile of user " + userProfile)
	return

}
