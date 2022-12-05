package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO Commentare la funzione
func (rt *_router) comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var textComment structs.TextComment

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")
	userID := ps.ByName("userID")

	// Check if the user I'm trying to follow has blocked me
	isBanned, err := rt.db.IsBanned(userID, ctx.Token)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to comment: userB has banned userA. userA: " + userID + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "You cannot unban a person which wasn't banned", http.StatusForbidden)
		return
	}

	isBanned, err = rt.db.IsBanned(ctx.Token, userID)

	if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		//^Aggiungere Forbidden come possibile risposta all'openapi
		rt.baseLogger.Error("Unable to comment: userA has banned userB. userA: " + userID + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&textComment)

	if err != nil {
		//^Aggiungere BadRequest come possibile risposta all'openapi
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Controlla se sto cercando di commentare una foto che non appartiene all'userID del path
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
	} else if owner != userID {
		// ^Unauthorized va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is trying to comment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = rt.db.Comment(photoID, ctx.Token, textComment.Text)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while commenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Photo commented")
}

// TODO Finire di commentare la funzione
func (rt *_router) unComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoOwner := ps.ByName("userID")
	photoID := ps.ByName("photoID")
	commentID := ps.ByName("commentID")

	// Controlla se sto cercando di commentare una foto che non appartiene all'userID del path
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

	err = rt.db.Uncomment(photoID, ctx.Token, commentID)

	if errors.Is(err, errorDefinition.ErrCommmentNotFound) {
		// ^BadRequest va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Comment not found")
		httpErrorResponse(rt, w, "Not Found, wrong request", http.StatusBadRequest)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while uncommenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo uncommented")

	//^Aggiungere StatusNoContent come possibile risposta all'openapi
	w.WriteHeader(http.StatusNoContent)

	return
}
