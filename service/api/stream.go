package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var following []string

	// Get page number from query string
	page := r.URL.Query().Get("page")

	// Convert page number to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		// If the page number is not a number, set it to 1
		pageInt = 0
	}

	if pageInt < 0 {
		pageInt = 0
	}

	following, err = rt.db.GetFollowing(ctx.Token)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("Error while getting following")
		httpErrorResponse(rt, w, "Internal Sever Error", http.StatusBadRequest)
		return
	}

	var photos []structs.Photo

	photos, err = rt.db.GetFollowingPhotosChrono(following, pageInt*100)

	if err == sql.ErrNoRows {
		// ^404 va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("No more photos to show")
		httpErrorResponse(rt, w, "NotFound", http.StatusNotFound)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("Error while getting photos")
		httpErrorResponse(rt, w, "Internal Sever Error", http.StatusBadRequest)
		return
	}

	// Send the response
	err = json.NewEncoder(w).Encode(photos)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("Error while encoding json")
		httpErrorResponse(rt, w, "Internal Sever Error", http.StatusBadRequest)
		return
	}
}
