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

	//Block/Unblock user endpoint
	//rt.router.PUT("/users/:userID/bans/:blockedID", rt.wrapAuth(rt.blockUser))
	rt.router.PUT("/users/:userID/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.blockUser)))
	//rt.router.DELETE("/users/:userID/bans/:blockedID", rt.wrapAuth(rt.blockUser))

	//Testing function
	rt.router.GET("/testing", rt.test)

	return rt.router
}
