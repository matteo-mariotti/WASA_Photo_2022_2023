package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")
	photoOwner := ps.ByName("userID")
	currentUser := ps.ByName("likeID")

	// Check if the user I'm trying to like a photo of a user who has blocked me
	isBanned, err := rt.db.IsBanned(photoOwner, ctx.Token)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + photoOwner + " had already banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to follow: userB has banned userA. userA: " + photoOwner + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unban a person which wasn't banned", http.StatusForbidden)
		return
	}

	// Controlla se sto cercando di mettere like ad una foto che non appartiene all'userID del path
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if err == sql.ErrNoRows {
		// ^Not Found va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	} else if owner != photoOwner {
		// ^StatusBadRequest va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is trying to comment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Controlla se il likeID coincide con l'user che è loggato
	if currentUser != ctx.Token {
		// ^Unauthorized va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is tying to impersonate someone else putting a like on a photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = rt.db.Like(photoID, ctx.Token)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while liking a photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Photo liked")
}

// TODO Finire di commentare la funzione
func (rt *_router) unlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoOwner := ps.ByName("userID")
	photoID := ps.ByName("photoID")
	currentUser := ps.ByName("likeID")

	// Controlla se sto cercando di togliere il like ad una foto che non appartiene all'userID del path
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if err == sql.ErrNoRows {
		// ^Not Found va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	} else if owner != photoOwner {
		// ^Unauthorized va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is trying to uncomment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Controlla se il likeID coincide con l'user che è loggato
	if currentUser != ctx.Token {
		// ^Unauthorized va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is tying to impersonate someone else putting a like on a photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = rt.db.Unlike(photoID, ctx.Token)

	if errors.Is(err, errorDefinition.ErrLikeNotFound) {
		// ^BadRequest va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Like not found")
		httpErrorResponse(rt, w, "Not Found, wrong request", http.StatusBadRequest)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while uncommenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo unlike")

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	return
}
