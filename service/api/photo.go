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

// TODO Commentare la funzione
func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// & Salva il file nella path /Photos
	var path = rt.photoPath

	// Inizia transazione
	// ^Internal Server Error va aggiunto all'openapi come possibile risposta
	err := rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Inserisci la riga nel db usando l'uuid come nome del nuovo file nella path ../../Photos
	newUuid, err := uuid.NewV4()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while generating UUID")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = rt.db.UploadPhoto(ps.ByName("userID"), newUuid.String())

	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while uploading photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Crea il file nella path ../../Photos
	f, err := os.Create(path + newUuid.String())
	defer f.Close()

	if err != nil {
		//Errore nella creazione del file, rollback
		rt.baseLogger.WithError(err).Error("Error while creating photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {

		//Errore nella copia del file
		rt.baseLogger.WithError(err).Error("Error while copying photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		//Rollback
		err = rt.db.Rollback()
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while rolling back")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		//Elimina il file creato
		err = os.Remove(path + newUuid.String())
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while removing photo from disk")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	rt.baseLogger.Info("Photo uploaded")
	rt.db.Commit()
	//Commit
}

// TODO Commentare la funzione
func (rt *_router) DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// & Salva il file nella path /Photos
	var path = rt.photoPath

	// Inizia transazione
	// ^Internal Server Error va aggiunto all'openapi come possibile risposta
	err := rt.db.StartTransaction()
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error while starting transaction")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Controlla se è il proprietario della foto
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	if err == sql.ErrNoRows {
		// ^Not Found va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		rt.db.Rollback()
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	} else if owner != ctx.Token {
		// ^Unauthorized va aggiunto all'openapi come possibile risposta
		rt.baseLogger.Error("User is trying to delete someone else's photo")
		httpErrorResponse(rt, w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rimuovi la riga dal db usando l'ID della foto
	fileIdentifier, err := rt.db.DeletePhoto((ps.ByName("photoID")))

	if errors.Is(err, errorDefinition.ErrPhotoNotFound) {
		// ^Not Found va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while deleting photo on DB")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	// Rimuovi il file dalla path /Photos
	err = os.Remove(path + fileIdentifier)

	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		//Errore nella eliminazione del file, rollback
		rt.baseLogger.WithError(err).Error("Error while deleting photo on disk")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo deleted")

	// Commit
	rt.db.Commit()

}

// TODO Commentare la funzione
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Check if the owner of the photo has not banned the user who is trying to access it
	owner, err := rt.db.GetPhotoOwner(ps.ByName("photoID"))
	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while getting photo owner")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if owner != ctx.Token {
		banned, err := rt.db.IsBanned(owner, ctx.Token)
		if err != nil {
			// ^Internal Server Error va aggiunto all'openapi come possibile risposta
			rt.baseLogger.WithError(err).Error("Error while checking if user is banned")
			httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if banned {
			// ^Forbidden va aggiunto all'openapi come possibile risposta
			rt.baseLogger.WithError(err).Error("User is banned")
			httpErrorResponse(rt, w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// Get the file from the path /Photos
	filename, err := rt.db.GetPhoto(ps.ByName("photoID"))

	if errors.Is(err, errorDefinition.ErrPhotoNotFound) {
		// ^Not Found va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Photo not found")
		httpErrorResponse(rt, w, "Not Found, wrong ID", http.StatusNotFound)
		return
	} else if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while getting photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Open the file
	file, err := os.Open(rt.photoPath + filename)
	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while opening photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//Copy the file to the response
	// Aggiorna l'api.yaml con il tipo image/*
	_, err = io.Copy(w, file)
	//Controllare se copy invia anche l'http status code
	if err != nil {
		// ^Internal Server Error va aggiunto all'openapi come possibile risposta
		rt.baseLogger.WithError(err).Error("Error while copying photo")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the action
	rt.baseLogger.Info("Photo sent")
	// Set the response code
	// ^ Forse è meglio un altro codice di risposta?
	// ! Perchè dice che è superfluo? Copy aggiunge anche lo stato http?
	//w.WriteHeader(http.StatusOK)
}

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

	// If not, get the profile
	likesResponse, err = rt.db.GetLikes(photoID, pageInt*30)

	if err == sql.ErrNoRows {
		//^Aggiungere NotFound come possibile risposta all'openapi
		rt.baseLogger.Error("No more likes are available")
		httpErrorResponse(rt, w, "No more likes", http.StatusNotFound)
		return
	} else if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the likes")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Prepare the JSON (it also sends the response to the client with the correct status code)
	err = json.NewEncoder(w).Encode(likesResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the likes of photo " + photoID)
	return
}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var likesResponse []structs.Comment

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

	// If not, get the profile
	likesResponse, err = rt.db.GetComments(photoID, pageInt*30)

	if err == sql.ErrNoRows {
		//^Aggiungere NotFound come possibile risposta all'openapi
		rt.baseLogger.Error("No more comments are available")
		httpErrorResponse(rt, w, "No more comments", http.StatusNotFound)
		return
	} else if err != nil {
		//^Aggiungere InternalServerError come possibile risposta all'openapi
		rt.baseLogger.Error("Error while getting the comments")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Prepare the JSON (it also sends the response to the client with the correct status code)
	err = json.NewEncoder(w).Encode(likesResponse)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		httpErrorResponse(rt, w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Log the action
	rt.baseLogger.Info("User " + ctx.Token + " has successfully got the comments of photo " + photoID)
	return

}
