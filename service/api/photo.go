package api

import (
	"WASA_Photo/service/api/reqcontext"
	"WASA_Photo/service/errorDefinition"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"

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
