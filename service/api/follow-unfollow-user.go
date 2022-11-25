package api

import (
	"WASA_Photo/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// UserA wants to follow userB

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("followerID")

	// Check if the user I'm trying to follow has blocked me
	isBanned, err := rt.db.IsBanned(userB, userA)

	// ^Internal Server Error va aggiunto all'openapi come possibile risposta
	if err != nil {
		errorLogger(rt, w, "Error while checking if user "+userB+" has blocked user "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ^Forbidden va aggiunto all'openapi come possibile risposta
	if isBanned {
		errorLogger(rt, w, "Unable to follow: userB has banned userA. userA: "+userA+" userB: "+userB, "You cannot follow a person which has banned you", http.StatusForbidden)
		return
	}

	// Check if I am trying to unblock myself
	if userA == userB {
		errorLogger(rt, w, "User is trying to follow himself: "+ctx.Token, "You cannot follow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(userA, userB)

	if err != nil {
		errorLogger(rt, w, "Error while checking if user "+userA+" is already following user "+userB, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isFollowing {
		errorLogger(rt, w, "UserA is already following userB. userA: "+userA+" userB: "+userB, "You cannot follow a person which you are already following", http.StatusBadRequest)
		return
	}

	err = rt.db.FollowUser(userA, userB)

	if err != nil {
		errorLogger(rt, w, "Error while following user "+userB+" from "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)
	rt.baseLogger.Info("UserA is now following userB. userA: " + userA + " userB: " + userB)
	return
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("followerID")

	// Check if I am trying to unblock myself
	if userA == userB {
		errorLogger(rt, w, "User is trying to unfollow himself: "+ctx.Token, "You cannot unfollow yourself", http.StatusBadRequest)
		return
	}

	// Check if userA is already following userB
	isFollowing, err := rt.db.IsFollowing(userA, userB)

	if err != nil {
		errorLogger(rt, w, "Error while checking if user "+userA+" is following user "+userB, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isFollowing {
		errorLogger(rt, w, "UserA is not following userB. userA: "+userA+" userB: "+userB, "You cannot unfollow a person which you are not following", http.StatusBadRequest)
		return
	}

	err = rt.db.UnfollowUser(userA, userB)

	if err != nil {
		errorLogger(rt, w, "Error while unfollowing user "+userB+" from "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)
	rt.baseLogger.Info("UserA has unfollowed userB. userA: " + userA + " userB: " + userB)
	return
}
