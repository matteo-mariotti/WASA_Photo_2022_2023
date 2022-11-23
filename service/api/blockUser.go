package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/api/utilities"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) blockUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	userA := ps.ByName("userID")
	userB := ps.ByName("blockedID")

	//Check if the user is using a valid token in the authentication header (if not, the request is blocked)
	if !ctx.Valid {
		rt.baseLogger.Error("Token is not valid: " + ctx.Token)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check if the user is acting on his own account
	if !utilities.CheckSelfAction(ctx.Token, userA) {
		rt.baseLogger.Error("User id trying to modify someone else's profile: " + ctx.Token)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO Implementare il ban di un utente

	// Need to check if I am trying to block myself
	if userA == userB {
		rt.baseLogger.Error("User id trying to block himself: " + ctx.Token)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Need to check if the userB is already blocked by userA
	// If not, block userB

	err := rt.db.BlockUser(userA, userB)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error while blocking userB from userA")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
