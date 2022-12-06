package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// getUserProfile parses the request extracting the user id, then, after checking the ban status of the user, it returns the user profile
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var profileResponse structs.Profile

	// Parsing the parameters from the request
	userProfile := ps.ByName("userID")

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

	// Check if the user I'm trying to look at has blocked me
	isBanned, err := rt.db.IsBanned(userProfile, ctx.Token)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userProfile + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to get the profile: userB has banned userA. (or viceversa) userB: " + userProfile + " userA: " + ctx.Token)
		httpErrorResponse(rt, w, "Fobidden", http.StatusForbidden)
		return
	}

	// Check if I have blocked the user I'm trying to look at
	isBannedViceversa, err := rt.db.IsBanned(ctx.Token, userProfile)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + ctx.Token + " has banned user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBannedViceversa {
		rt.baseLogger.Error("Unable to get the profile: userA has banned userB. (or viceversa) userB: " + ctx.Token + " userA: " + userProfile)
		httpErrorResponse(rt, w, "Fobidden", http.StatusForbidden)
		return
	}

	// If not, get the profile
	username, err := rt.db.GetName(userProfile)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Error("This user does not exist!" + userProfile)
		httpErrorResponse(rt, w, "UserID is not valid", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while getting the username of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the username
	profileResponse.Username = username

	// Get the number of followers
	follower, err := rt.db.GetFollowerNumber(userProfile)

	if err != nil {
		rt.baseLogger.Error("Error while getting the number of followers of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of followers
	profileResponse.Follower = follower

	// Get the number of following
	following, err := rt.db.GetFollowingNumber(userProfile)

	if err != nil {
		rt.baseLogger.Error("Error while getting the number of following of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of following
	profileResponse.Following = following

	// Get the number of photos
	photoNumber, err := rt.db.GetPhotosNumber(userProfile)

	if err != nil {
		rt.baseLogger.Error("Error while getting the number of photos of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the number of photos
	profileResponse.PhotoCount = photoNumber

	// Get the photos (using the offset)
	photos, err := rt.db.GetPhotos(userProfile, pageInt*30)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Error("No more photos are available in this user profile:  " + userProfile)
		httpErrorResponse(rt, w, "404 Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while getting the photos of user " + userProfile)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the photos
	profileResponse.Photo = photos

	//Prepare the JSON (it also sends the response to the client with the correct status code)
	err = json.NewEncoder(w).Encode(profileResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the profile of user " + userProfile)
	return

}
