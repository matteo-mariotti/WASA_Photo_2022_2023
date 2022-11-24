package api

import (
	"WASA_Photo/service/api/reqcontext"
	"database/sql"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// wrapAuth parses the request and checks if the bearer token in the authorization header is valid or not
func (rt *_router) wrapAuth(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		//Check if the token is valid
		token := r.Header.Get("Authorization")
		splitToken := strings.Split(token, "Bearer ")
		token = splitToken[1]
		result, err := rt.db.ValidToken(token)

		if err != nil && err != sql.ErrNoRows {
			ctx.Logger.WithError(err).Error("Error while checking if the token is valid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Save the result in the context
		ctx.Valid = result
		ctx.Token = token

		if result {
			rt.baseLogger.Info("User is using a valid token: ", ctx.Token)
		} else {
			rt.baseLogger.Warning("User is using an invalid token: ", ctx.Token)
		}

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
