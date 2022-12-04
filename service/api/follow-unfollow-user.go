package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// followerID wants to follow useID

	// Parsing the parameters from the request
	userID := ps.ByName("userID")
	followerID := ps.ByName("followerID")

	if followerID != ctx.Token {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("User id trying to fact as someone else: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot impersonate someone else", http.StatusForbidden)
		return
	}

	// Check if the user I'm trying to follow has blocked me
	isBanned, err := rt.db.IsBanned(userID, followerID)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to follow: userID has banned followerID. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "You cannot follow a person which has banned you", http.StatusForbidden)
		return
	}

	// Check if I am trying to unblock myself
	if userID == followerID {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("User is trying to follow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(followerID, userID)

	if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + followerID + " is following user " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isFollowing {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("FollowerID is already following userID. FollowerID: " + followerID + " userID: " + userID)
		httpErrorResponse(rt, w, "You cannot follow a person which you are already following", http.StatusBadRequest)
		return
	}

	err = rt.db.FollowUser(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserID or followerID not found. UserID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or FollowerID not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.WithError(err).Error("Error while following user. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("FollowerID is now following userID. FollowerID: " + followerID + " userID: " + userID)
	return
}

// TODO Commentare la funzione

// TODO RISCRIERE LA FUNZIONE CON LE VARIABILI CHIAMATE DIVERSAMENTE PER NON CONFONDERLE CON LE ALTRE
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parsing the parameters from the request
	userID := ps.ByName("userID")
	followerID := ps.ByName("followerID")

	if followerID != ctx.Token {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("User id trying to act as someone else: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot impersonate someone else", http.StatusForbidden)
		return
	}

	// Check if I am trying to unblock myself
	if userID == followerID {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("User is trying to unfollow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unfollow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserID or followerID not found. UserID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or followerID not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if " + followerID + " is following " + userID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isFollowing {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("FollowerID is not following userID. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "You cannot unfollow a person which you are not following", http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(followerID, userID)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserID or followerID not found. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "UserID or followerID not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while unfollowing user. userID: " + userID + " followerID: " + followerID)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("FollowerID has unfollowed userID. userID: " + userID + " followerID: " + followerID)
	return
}
