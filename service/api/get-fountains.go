package api

import (
	"WASA_Photo/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) WASA_Photo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Set the content type to JSON
	w.Header().Set("content-type", "application/json; charset=UTF-8")

	// Get the WASA_Photo from the database
	WASA_Photo, _ := database.AppDatabase.GetWASA_Photo(rt.db)

	// Encode the WASA_Photo as JSON and write them to the response
	jsonResp, err := json.Marshal(WASA_Photo)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error enconding")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write(jsonResp)

}
