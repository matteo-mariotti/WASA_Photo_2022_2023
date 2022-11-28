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

	// UserA wants to follow userB

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("followerID")

	// Check if the user I'm trying to follow has blocked me
	isBanned, err := rt.db.IsBanned(userB, userA)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userA + " had already banned user " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to follow: userB has banned userA. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "You cannot unban a person which wasn't banned", http.StatusForbidden)
		return
	}

	// Check if I am trying to unblock myself
	if userA == userB {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("User is trying to follow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(userA, userB)

	if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userA + " is following user " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isFollowing {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserA is already following userB. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "You cannot follow a person which you are already following", http.StatusBadRequest)
		return
	}

	err = rt.db.FollowUser(userA, userB)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserA or userB not found. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "UserA or userB not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while following user. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("UserA is now following userB. userA: " + userA + " userB: " + userB)
	return
}

// TODO Commentare la funzione
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("followerID")

	// Check if I am trying to unblock myself
	if userA == userB {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("User is trying to unfollow himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unfollow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(userA, userB)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserA or userB not found. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "UserA or userB not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userA + " is following user " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isFollowing {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserA is not following userB. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "You cannot unfollow a person which you are not following", http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(userA, userB)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		//^Aggiungere StatusBadRequest come possibile risposta all'openapi
		rt.baseLogger.Error("UserA or userB not found. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "UserA or userB not found", http.StatusBadRequest)
		return
	} else if err != nil {
		//^Aggiungere StatusInternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while unfollowing user. userA: " + userA + " userB: " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("UserA has unfollowed userB. userA: " + userA + " userB: " + userB)
	return
}
