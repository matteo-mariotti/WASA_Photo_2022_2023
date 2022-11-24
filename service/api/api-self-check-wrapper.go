package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/api/utilities"
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// wrapSelf parses the request (the userID field) and checks if the user is acting on his own account
func (rt *_router) wrapSelf(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		userA := ps.ByName("userID")

		// Check if the user is acting on his own account
		if !utilities.CheckSelfAction(ctx.Token, userA) {
			// Log the error
			rt.baseLogger.Error("User id trying to modify someone else's profile: " + ctx.Token)

			// Return the error
			w.Header().Set("Content-Type", "application/json") // Set content to JSON
			// ?     	Non sarebbe meglio ritornare un errore Forbidden?
			w.WriteHeader(http.StatusUnauthorized)                                                               // Set status to unauthorized
			json.NewEncoder(w).Encode(structs.GenericError{Message: "You cannot modify someone else's profile"}) // Encode the error and send it
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
