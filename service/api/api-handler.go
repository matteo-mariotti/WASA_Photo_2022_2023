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
	rt.router.PUT("/users/:userID/followers/:followerID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.followUser))))
	rt.router.DELETE("/users/:userID/followers/:followerID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.unfollowUser))))

	//Upload/Delete photo route
	rt.router.POST("/users/:userID/photos", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.UploadPhoto))))
	rt.router.DELETE("/users/:userID/photos/:photoID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.DeletePhoto))))

	//Add/Delete comment route
	rt.router.POST("/users/:userID/photos/:photoID/comments", rt.wrap(rt.wrapAuth(rt.comment)))
	rt.router.DELETE("/users/:userID/photos/:photoID/comments/:commentID", rt.wrap(rt.wrapAuth(rt.unComment)))

	//Get photo route
	rt.router.GET("/photos/:photoID", rt.wrap(rt.getPhoto))

	// Add/Delete like route
	rt.router.PUT("/users/:userID/photos/:photoID/likes/:likeID", rt.wrap(rt.wrapAuth(rt.like)))
	rt.router.DELETE("/users/:userID/photos/:photoID/likes/:likeID", rt.wrap(rt.wrapAuth(rt.unlike)))

	//Testing function
	rt.router.GET("/testing", rt.test)

	return rt.router
}
