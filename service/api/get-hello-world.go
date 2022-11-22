package api

import (
	"WASA_Photo/service/database"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	version, _ := database.AppDatabase.GetVersion(rt.db)

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
	_, _ = w.Write([]byte("\nDatabase version: " + version))

}
