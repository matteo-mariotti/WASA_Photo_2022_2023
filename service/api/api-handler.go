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

	//Session routes
	rt.router.POST("/session", rt.wrap(rt.login))

	//Block/Unblock routes
	rt.router.PUT("/users/:userID/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.banUser))))
	rt.router.DELETE("/users/:userID/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.unbanUser))))

	//Follow/Unfollow route
	rt.router.PUT("/users/:userID/followers/:followerID", rt.wrap(rt.wrapAuth(rt.followUser)))
	rt.router.DELETE("/users/:userID/followers/:followerID", rt.wrap(rt.wrapAuth(rt.unfollowUser)))

	//Testing function
	rt.router.GET("/testing", rt.test)

	return rt.router
}
