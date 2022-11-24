package api

import (
	"WASA_Photo/service/structs"
	"encoding/json"
	"net/http"
)

func httpCheckError(rt *_router, w http.ResponseWriter, message string, err error, httpError int) bool {

	if err != nil {
		// Log the error
		rt.baseLogger.Error(message + err.Error())

		// Return the error
		w.Header().Set("Content-Type", "application/json") // Set content to JSON

		w.WriteHeader(httpError)                                          // Set status to the correct error
		json.NewEncoder(w).Encode(structs.GenericError{Message: message}) // Encode the error and send it
		return true
	}
	return false
}

func errorLogger(rt *_router, w http.ResponseWriter, logMessage string, responseMessage string, httpError int) {
	// Log the error
	rt.baseLogger.Error(logMessage)

	// Return the error
	w.Header().Set("Content-Type", "application/json")                        // Set content to JSON
	w.WriteHeader(httpError)                                                  // Set status to the correct error
	json.NewEncoder(w).Encode(structs.GenericError{Message: responseMessage}) // Encode the error and send it
	return
}
