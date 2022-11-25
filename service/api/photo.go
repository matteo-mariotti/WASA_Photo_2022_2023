package api

import (
	"WASA_Photo/service/api/reqcontext"
	"database/sql"
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
		errorLogger(rt, w, "Error while starting transaction", "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Inserisci la riga nel db usando l'uuid come nome del nuovo file nella path ../../Photos
	newUuid, err := uuid.NewV4()
	if err != nil {
		errorLogger(rt, w, "Error while generating uuid", "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rt.db.UploadPhoto(ps.ByName("userID"), newUuid.String())

	f, err := os.Create(path + newUuid.String())
	defer f.Close()

	if err != nil {
		//Errore nella creazione del file, rollback
		rt.baseLogger.WithError(err).Error("Errore")
		errorLogger(rt, w, "Error while creating file", "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	_, err = io.Copy(f, r.Body)
	if err != nil {
		//Errore nella copia del file
		errorLogger(rt, w, "Error while copying file", "Internal Server Error", http.StatusInternalServerError)
		//Elimina il file creato
		os.Remove(path + newUuid.String())
		//Rollback
		rt.db.Rollback()
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
		errorLogger(rt, w, "Error while starting transaction", "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Controlla se è il proprietario della foto
	owner, err := rt.db.GetPhotoOwner((ps.ByName("photoID")))

	// ^Not Found va aggiunto all'openapi come possibile risposta
	if err == sql.ErrNoRows {
		errorLogger(rt, w, "Photo not found, user is using the wrong ID", "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		errorLogger(rt, w, "Error while getting photo owner", "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	// ^Unauthorized va aggiunto all'openapi come possibile risposta
	if owner != ctx.Token {
		errorLogger(rt, w, "User is trying to delete someone else's photo", "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rimuovi la riga dal db usando l'ID della foto
	fileIdentifier, err := rt.db.DeletePhoto((ps.ByName("photoID")))

	// ^Not Found va aggiunto all'openapi come possibile risposta
	if err == sql.ErrNoRows {
		errorLogger(rt, w, "Photo not found, user is using the wrong ID", "Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		//Errore nella rimozione del file, rollback
		rt.baseLogger.WithError(err).Error("Errore DB")
		errorLogger(rt, w, "Error while deleting file", "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	// Rimuovi il file dalla path /Photos
	err = os.Remove(path + fileIdentifier)

	if err != nil {
		//Errore nella eliminazione del file, rollback
		rt.baseLogger.WithError(err).Error("Errore")
		errorLogger(rt, w, "Error while deleting file", "Internal Server Error", http.StatusInternalServerError)
		rt.db.Rollback()
		return
	}

	rt.baseLogger.Info("Photo deleted")
	rt.db.Commit()
	//Commit

}
