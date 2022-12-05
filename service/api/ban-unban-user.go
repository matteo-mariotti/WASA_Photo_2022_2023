package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// banUser parses the request extracting the user id to ban and the user id of the user who is banning the user then,
// after doing some checks, it bans the user
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("blockedID")

	// Check if I am trying to ban myself
	if userA == userB {
		rt.baseLogger.Error("User is trying to ban himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot ban yourself", http.StatusBadRequest)
		return
	}

	// Check if the userB is already banned by userA
	isBanned, err := rt.db.IsBanned(userA, userB)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userA + " has already banned user " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("User " + userA + " had already banned user " + userB)
		httpErrorResponse(rt, w, "You cannot ban a person which was already banned", http.StatusBadRequest)
		return
	}

	// If not, ban userB
	err = rt.db.BanUser(userA, userB)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		rt.baseLogger.Error("User " + userA + " or " + userB + " not found")
		httpErrorResponse(rt, w, "UserA or UserB not found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while banning user " + userB + " from user " + userA)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)

	//Log the action
	rt.baseLogger.Info("UserB banned from userA. userA: " + userA + " userB: " + userB)
	return

}

// unbanUser parses the request extracting the user id to unban and the user id of the user who is unbaning the user then,
// after doing some checks, it unbans the user
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("blockedID")

	// Check if I am trying to unblock myself
	if userA == userB {
		// ^StatusBadRequest va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is trying to unban himself: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unban yourself", http.StatusBadRequest)
		return
	}

	// Check if the userB was blocked by userA
	isBanned, err := rt.db.IsBanned(userA, userB)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userA + " had already banned user " + userB)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if !isBanned {
		rt.baseLogger.Error("User " + userA + " had not banned user " + userB)
		httpErrorResponse(rt, w, "You cannot unban a person which wasn't banned", http.StatusBadRequest)
		return
	}

	// If was banned, unban userB
	err = rt.db.UnbanUser(userA, userB)

	if errors.Is(err, errorDefinition.ErrUserNotFound) {
		rt.baseLogger.Error("User " + userA + " or " + userB + " not found")
		httpErrorResponse(rt, w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while unbanning user " + userB + " from user " + userA)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If the operation was successful return the status code 204 No Content
	w.WriteHeader(http.StatusNoContent)

	//Log the action
	rt.baseLogger.Info("UserB unbanned from userA. userA: " + userA + " userB: " + userB)
	return
}
