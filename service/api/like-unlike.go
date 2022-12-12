package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// like is the function that handles the request to like a photo
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")
	photoOwner := ps.ByName("userID")
	currentUser := ps.ByName("likeID")

	// Check if the user I'm trying to like a photo of a user who has blocked me
	isBanned, err := rt.db.IsBanned(photoOwner, ctx.Token)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + photoOwner + " had already banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to follow: userB has banned userA. userA: " + photoOwner + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if the user I'm trying to like a photo of a user who has blocked me
	isBanned, err = rt.db.IsBanned(photoOwner, ctx.Token)

	if err != nil {
		rt.baseLogger.Error("Error while checking if user " + photoOwner + " had already banned user " + ctx.Token)
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("Unable to follow: userA has banned userB. userA: " + photoOwner + " userB: " + ctx.Token)
		httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if I'm trying to like a photo that doesn't belong to the userID in the path
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	} else if owner != photoOwner {
		rt.baseLogger.Error("User is trying to comment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if the user is acting as himself
	if currentUser != ctx.Token {
		rt.baseLogger.Error("User is tying to impersonate someone else putting a like on a photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Like the photo
	err = rt.db.Like(photoID, ctx.Token)

	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while liking a photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If everything went well, return a 204
	w.WriteHeader(http.StatusNoContent)

	// Log the action
	rt.baseLogger.Info("Photo liked")
}

// unlike is the function that handles the request to unlike a photo
func (rt *_router) unlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parsing the parameters from the request
	photoOwner := ps.ByName("userID")
	photoID := ps.ByName("photoID")
	currentUser := ps.ByName("likeID")

	// Check if the user I'm trying to unlike a photo is the same as the one in the path
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	} else if owner != photoOwner {
		rt.baseLogger.Error("User is trying to uncomment a photo that doesn't belong to the user in the path")
		httpErrorResponse(rt, w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if the user is acting as himself
	if currentUser != ctx.Token {
		rt.baseLogger.Error("User is tying to impersonate someone else putting a like on a photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Unlike the photo
	err = rt.db.Unlike(photoID, ctx.Token)

	if errors.Is(err, errorDefinition.ErrLikeNotFound) {
		rt.baseLogger.WithError(err).Error("Like not found")
		httpErrorResponse(rt, w, "Not Found, wrong request", http.StatusBadRequest)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while uncommenting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo unlike")

	// If everything went well, return a 204
	w.WriteHeader(http.StatusNoContent)
}
