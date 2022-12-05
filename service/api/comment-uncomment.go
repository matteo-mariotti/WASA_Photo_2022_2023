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

// comment parses the request extracting the user id and the new username, it then checks if the owner of the photo has banned me or viceversa, if not,
// it checks wheter the photo owner corresponds to the one in the path, at the end, it comments the photo
func (rt *_router) comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var textComment structs.TextComment

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")
	userID := ps.ByName("userID")

	// Check if the user I'm trying to comment the photo of a user that has blocked me
	isBanned, err := rt.db.IsBanned(userID, ctx.Token)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to comment: userB has banned userA. userA: " + userID + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "This user has blocked you ", http.StatusConflict)
		return
	}

	// Check if I have blocked the user I'm trying to comment the photo
	isBanned, err = rt.db.IsBanned(ctx.Token, userID)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + userID + " has banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to comment: userA has banned userB. userA: " + userID + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "You blocked this user", http.StatusForbidden)
		return
	}

	// Get the text of the comment
	err = json.NewDecoder(r.Body).Decode(&textComment)

	if err != nil {
		rt.baseLogger.Error("Error while decoding the comment")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if I'm trying to comment a photo that doesn't belong to the user in the path
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if err == sql.ErrNoRows {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	} else if owner != userID {
		rt.baseLogger.Error("User is trying to comment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	// Comment the photo
	err = rt.db.Comment(photoID, ctx.Token, textComment.Text)

	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while commenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Photo commented")
}

// uncomment parses the request extracting the user id and the new username, then,
// it checks wheter the photo owner corresponds to the one in the path, at the end, it comments the photo
func (rt *_router) unComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoOwner := ps.ByName("userID")
	photoID := ps.ByName("photoID")
	commentID := ps.ByName("commentID")

	// Controlla se sto cercando di commentare una foto che non appartiene all'userID del path
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if err == sql.ErrNoRows {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	} else if owner != photoOwner {
		rt.baseLogger.Error("User is trying to uncomment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	// UnComment the photo
	err = rt.db.Uncomment(photoID, ctx.Token, commentID)

	if errors.Is(err, errorDefinition.ErrCommmentNotFound) {
		rt.baseLogger.WithError(err).Error("Comment not found")
		httpErrorResponse(rt, w, "Not Found, wrong request", http.StatusBadRequest)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while uncommenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo uncommented")

	w.WriteHeader(http.StatusNoContent)

	return
}
