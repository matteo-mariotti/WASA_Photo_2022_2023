package api

import (
	"WASA_Photo/service/api/reqcontext"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrapAuth(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

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
