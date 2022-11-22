package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	//Session endpoint
	rt.router.POST("/session", rt.wrap(rt.login))

	rt.router.PUT("/users/.userID/bans/.blockedID", rt.wrap(rt.login))

	rt.router.GET("/WASA_Photo", rt.WASA_Photo)
	rt.router.POST("/WASA_Photo", rt.create_fountain)

	return rt.router
}
