package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"WASA_Photo/service/structs"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// uploadPhoto uploads a photo to the DB and to the disk
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var path = rt.photoPath

	// Start transaction
	err := rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert the photo in the DB
	newUuid, err := uuid.NewV4()
	if err != nil {
		err1 := rt.db.Rollback()
		if err1 != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		rt.baseLogger.WithError(err).Error("Error while generating UUID")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := rt.db.GetToken(ps.ByName("username"))
	if err != nil {
		// Error, rollback
		rt.baseLogger.WithError(err).Error("Error while getting token")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	err = rt.db.UploadPhoto(token, newUuid.String())

	if err != nil {
		err1 := rt.db.Rollback()
		if err1 != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		rt.baseLogger.WithError(err).Error("Error while uploading photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the physical file in  ../../Photos
	f, err := os.Create(path + newUuid.String())

	if err != nil {
		// Error, rollback
		rt.baseLogger.WithError(err).Error("Error while creating photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	defer f.Close()

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		rt.baseLogger.Info("Form cose 1")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		rt.baseLogger.Info("Form cose 2")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the photo in the new file
	_, err = io.Copy(f, file)
	if err != nil {

		// Error, rollback
		rt.baseLogger.WithError(err).Error("Error while copying photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Delete the file
		err = os.Remove(path + newUuid.String())
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while removing photo from disk")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Commit
	err = rt.db.Commit()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while committing")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rt.baseLogger.Info("Photo uploaded")

}

// DeletePhoto deletes a photo from the DB and from the disk
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var path = rt.photoPath

	// Start transaction
	err := rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Controlla se è il proprietario della foto
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
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
	} else if owner != ctx.Token {
		rt.baseLogger.Error("User is trying to delete someone else's photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Remove the photo from the DB
	fileIdentifier, err := rt.db.DeletePhoto((ps.ByName("photoID")))

	if errors.Is(err, errorDefinition.ErrPhotoNotFound) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while deleting photo on DB")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Remove the file from the disk
	err = os.Remove(path + fileIdentifier)

	if err != nil {
		// Error, rollback
		rt.baseLogger.WithError(err).Error("Error while deleting photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Unable to rollback")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo deleted")

	// Commit
	err = rt.db.Commit()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Unable to rollback")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
	}

}

// getPhoto returns the photo with the given ID
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Check if the owner of the photo has not banned the user who is trying to access it
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if owner != ctx.Token {
		banned, err := rt.db.IsBanned(owner, ctx.Token)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}
		banned, err = rt.db.IsBanned(ctx.Token, owner)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}

	}

	// Get the file from the path /Photos
	filename, err := rt.db.GetPhoto(ps.ByName("photoID"))

	if errors.Is(err, errorDefinition.ErrPhotoNotFound) {
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Open the file
	file, err := os.Open(rt.photoPath + filename)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while opening photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//http.ServeFile(w, r, rt.photoPath+filename)

	// Copy the file to the response
	_, err = io.Copy(w, file)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while copying photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo sent")
}

// getLikes returns the likes of the photo with the given ID
func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var likesResponse []structs.Like

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")

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

	// Check if the owner of the photo has not banned the user who is trying to access it
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if owner != ctx.Token {
		banned, err := rt.db.IsBanned(owner, ctx.Token)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}
		banned, err = rt.db.IsBanned(ctx.Token, owner)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// If not, get the likes
	likesResponse, err = rt.db.GetLikes(photoID, pageInt*30, ctx.Token)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Error("No more likes are available")
		httpErrorResponse(rt, w, "No more likes", http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while getting the likes")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the JSON (it also sends the response to the client with the correct status code)
	err = json.NewEncoder(w).Encode(likesResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the likes of photo " + photoID)
}

// getComments returns the comments of the photo with the given ID
func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var commentsResponse []structs.Comment

	// Parsing the parameters from the request
	photoID := ps.ByName("photoID")

	// Check if the owner of the photo has not banned the user who is trying to access it
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if owner != ctx.Token {
		banned, err := rt.db.IsBanned(owner, ctx.Token)

		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}

		banned, err = rt.db.IsBanned(ctx.Token, owner)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}

	}

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

	// If not, get the profile
	commentsResponse, err = rt.db.GetComments(photoID, pageInt*30, ctx.Token)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Info("No more comments are available")
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		rt.baseLogger.Error("Error while getting the comments")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the JSON (it also sends the response to the client with the correct status code)
	err = json.NewEncoder(w).Encode(commentsResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the comments of photo " + photoID)
}
