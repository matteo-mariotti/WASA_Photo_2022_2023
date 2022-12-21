package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Session routes
	rt.router.POST("/session", rt.wrap(rt.login))

	// Block/Unblock routes
	rt.router.PUT("/users/:username/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.banUser))))
	rt.router.DELETE("/users/:username/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.unbanUser))))
	rt.router.GET("/users/:username/bans/:blockedID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.banStatus))))

	// Follow/Unfollow route
	rt.router.PUT("/users/:username/followers/:followerID", rt.wrap(rt.wrapAuth(rt.followUser)))
	rt.router.DELETE("/users/:username/followers/:followerID", rt.wrap(rt.wrapAuth(rt.unfollowUser)))
	rt.router.GET("/users/:username/followers/:followerID", rt.wrap(rt.wrapAuth(rt.followStatus)))

	// Upload/Delete photo route
	rt.router.POST("/users/:username/photos", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.uploadPhoto))))
	rt.router.DELETE("/users/:username/photos/:photoID", rt.wrap(rt.wrapAuth(rt.wrapSelf(rt.deletePhoto))))

	// Add/Delete comment route
	rt.router.POST("/users/:username/photos/:photoID/comments", rt.wrap(rt.wrapAuth(rt.comment)))
	rt.router.DELETE("/users/:username/photos/:photoID/comments/:commentID", rt.wrap(rt.wrapAuth(rt.unComment)))

	// Get photo route
	rt.router.GET("/photos/:photoID", rt.wrap(rt.wrapAuth(rt.getPhoto)))
	rt.router.GET("/photos/:photoID/likes", rt.wrap(rt.wrapAuth(rt.getLikes)))
	rt.router.GET("/photos/:photoID/comments", rt.wrap(rt.wrapAuth(rt.getComments)))

	// Add/Delete like route
	rt.router.PUT("/users/:username/photos/:photoID/likes/:likeID", rt.wrap(rt.wrapAuth(rt.like)))
	rt.router.DELETE("/users/:username/photos/:photoID/likes/:likeID", rt.wrap(rt.wrapAuth(rt.unlike)))

	// Change username route
	rt.router.PUT("/users/:username/username", rt.wrap(rt.wrapAuth(rt.changeUsername)))

	// Get user profile route
	rt.router.GET("/users/:username", rt.wrap(rt.wrapAuth(rt.getUserProfile)))
	rt.router.GET("/users/:username/photos", rt.wrap(rt.wrapAuth(rt.getPhotos)))

	// Get users filtered list
	rt.router.GET("/users", rt.wrap(rt.wrapAuth(rt.getUsers)))

	// Stream route
	rt.router.GET("/stream", rt.wrap(rt.wrapAuth(rt.stream)))

	return rt.router
}
