package api

import (
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) create_fountain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var fountain structs.Fountain

	err := json.NewDecoder(r.Body).Decode(&fountain)
	_ = r.Body.Close()
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Wrong JSON received")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !fountain.IsValidFountain() {
		rt.baseLogger.WithError(err).Warning("Fountain is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Insert the fountain in the database
	err = rt.db.InsertFountain(fountain)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error inserting fountain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}
