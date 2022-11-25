package api

import (
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"
)

func errorLogger(rt *_router, w http.ResponseWriter, logMessage string, responseMessage string, httpError int) {
	// Log the error
	rt.baseLogger.Error(logMessage)

	// Return the error
	w.Header().Set("Content-Type", "application/json")                        // Set content to JSON
	w.WriteHeader(httpError)                                                  // Set status to the correct error
	json.NewEncoder(w).Encode(structs.GenericError{Message: responseMessage}) // Encode the error and send it
	return
}
