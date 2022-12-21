package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"WASA_Photo/service/structs"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// followUser is the function that handles the follow user request
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Note: followerID wants to follow username

	// Parsing the parameters from the request
	userID := ps.ByName("username")
	followerID := ps.ByName("followerID")

	userID, err := rt.db.GetToken(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	followerID, err = rt.db.GetToken(followerID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Check that the user is not acting as someone else
	if followerID != ctx.Token {
		rt.baseLogger.Error("User id trying to fact as someone else: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot impersonate someone else", http.StatusForbidden)
		return
	}

	// Check if the user I'm trying to follow has blocked me
	isBanned, err := rt.db.IsBanned(userID, followerID)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to follow: userID has banned followerID. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	// Checking the viceversa
	isBanned, err = rt.db.IsBanned(followerID, userID)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + followerID + " has banned user " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to follow: followerID has banned userID. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if I am trying to follow myself
	if userID == followerID {
		rt.baseLogger.Error("User is trying to follow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(followerID, userID)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + followerID + " is following user " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isFollowing {
		rt.baseLogger.Error("FollowerID is already following userID. FollowerID: " + followerID + " userID: " + userID)
		httpErrorResponse(rt, w, "You cannot follow a person which you are already following", http.StatusConflict)
		return
	}

	// Follow the user
	err = rt.db.FollowUser(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		rt.baseLogger.Error("UserID or followerID not found. UserID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or FollowerID not found", http.StatusBadRequest)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while following user. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("FollowerID is now following userID. FollowerID: " + followerID + " userID: " + userID)
}

// unfollowUser is the function that handles the unfollow user request
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userID := ps.ByName("username")
	followerID := ps.ByName("followerID")

	userID, err := rt.db.GetToken(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	followerID, err = rt.db.GetToken(followerID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Check that the user is not acting as someone else
	if followerID != ctx.Token {
		rt.baseLogger.Error("User id trying to act as someone else: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot impersonate someone else", http.StatusForbidden)
		return
	}

	// Check if I am trying to unfollow myself
	if userID == followerID {
		rt.baseLogger.Error("User is trying to unfollow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unfollow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is following userB
	isFollowing, err := rt.db.IsFollowing(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		rt.baseLogger.Error("UserID or followerID not found. UserID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or followerID not found", http.StatusBadRequest)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while checking if " + followerID + " is following " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isFollowing {
		rt.baseLogger.Error("FollowerID is not following userID. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "You cannot unfollow a person which you are not following", http.StatusConflict)
		return
	}

	// Unfollow the user
	err = rt.db.UnfollowUser(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		rt.baseLogger.Error("UserID or followerID not found. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or followerID not found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while unfollowing user. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("FollowerID has unfollowed userID. userID: " + userID + " followerID: " + followerID)
}

func (rt *_router) followStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userID := ps.ByName("username")
	followerID := ps.ByName("followerID")

	userID, err := rt.db.GetToken(userID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	followerID, err = rt.db.GetToken(followerID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(followerID, userID)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + followerID + " is following user " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isFollowing {
		rt.baseLogger.Error("FollowerID is already following userID. FollowerID: " + followerID + " userID: " + userID)
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(structs.Status{StatusRes: true}) // Encode the error and send it
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error encoding error response")
			w.WriteHeader(http.StatusInternalServerError) // Set status to the correct error
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(structs.Status{StatusRes: false}) // Encode the error and send it
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding error response")
		w.WriteHeader(http.StatusInternalServerError) // Set status to the correct error
	}
	return

}
