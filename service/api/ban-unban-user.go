package api

import (
	"WASA_Photo/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("blockedID")

	// Check if I am trying to ban myself
	if userA == userB {
		errorLogger(rt, w, "User is trying to ban himself: "+ctx.Token, "You cannot ban yourself", http.StatusBadRequest)
		return
	}

	// Check if the userB is already banned by userA
	res, err := rt.db.IsBanned(userA, userB)

	// ^Internal Server Error va aggiunto all'openapi come possibile risposta
	if err != nil {
		errorLogger(rt, w, "Error while checking if user "+userA+" has already banned user "+userB, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if res {
		errorLogger(rt, w, "UserB was already banned by userA. userA: "+userA+" userB: "+userB, "You cannot ban a person which was already banned", http.StatusBadRequest)
		return
	}

	// If not, ban userB
	err = rt.db.BanUser(userA, userB)

	//^Aggiungere InternalServerError come possibile risposta all'openapi
	if err != nil {
		errorLogger(rt, w, "Error while banning user "+userB+" from "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return 204
	w.WriteHeader(http.StatusNoContent)

	//Log the action
	rt.baseLogger.Info("UserB banned from userA. userA: " + userA + " userB: " + userB)
	return

}

// unbanUser parses the request extracting the user id to unban and the user id of the user who is unbaning the user then, after doing some checks, it unban the user
// TODO Finire di commentare con cosa ritorna
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("blockedID")

	// Check if I am trying to unblock myself
	if userA == userB {
		errorLogger(rt, w, "User is trying to unban himself: "+ctx.Token, "You cannot unban yourself", http.StatusBadRequest)
		return
	}

	// Check if the userB was blocked by userA
	res, err := rt.db.IsBanned(userA, userB)

	//^Aggiungere InternalServerError come possibile risposta all'openapi
	if err != nil {
		errorLogger(rt, w, "An error occured while checking if userB "+userB+" was banned by userA "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !res {
		errorLogger(rt, w, "UserB was not banned by userA. userA: "+userA+" userB: "+userB, "You cannot unban a person which wasn't banned", http.StatusBadRequest)
		return
	}

	// If was banned, unban userB
	err = rt.db.UnbanUser(userA, userB)

	//^Aggiungere InternalServerError come possibile risposta all'openapi
	if err != nil {
		errorLogger(rt, w, "An error occured while unbanning userB "+userB+" from: "+userA, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// If the operation was successful return the status code 204 No Content
	w.WriteHeader(http.StatusNoContent)
	rt.baseLogger.Info("UserB unbanned from userA. userA: " + userA + " userB: " + userB)
	return
}
