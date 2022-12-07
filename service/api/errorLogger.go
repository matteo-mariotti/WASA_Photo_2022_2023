package api

import (
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"
)

func httpErrorResponse(rt *_router, w http.ResponseWriter, responseMessage string, httpError int) {

	// Return the error
	w.Header().Set("Content-Type", "application/json")                               // Set content to JSON
	w.WriteHeader(httpError)                                                         // Set status to the correct error
	err := json.NewEncoder(w).Encode(structs.GenericError{Message: responseMessage}) // Encode the error and send it
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding error response")
		w.WriteHeader(http.StatusInternalServerError) // Set status to the correct error

	}
}
