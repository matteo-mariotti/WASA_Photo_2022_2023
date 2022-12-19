package api

import (
	"WASA_Photo/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// wrapSelf parses the request (the userID field) and checks if the user is acting on his own account in path like /users/:userID/... where userID is the ID of the user that should be acting on his own account
func (rt *_router) wrapSelf(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		userID := ps.ByName("username")

		userID, err := rt.db.GetToken(userID)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error getting user token")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if the user is acting on his own account
		if ctx.Token != userID {
			// Log the error
			rt.baseLogger.Error("User id trying to modify someone else's profile: " + ctx.Token)

			// Return the error
			httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
