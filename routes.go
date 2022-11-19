package main

import "github.com/gin-gonic/gin"

func addAPIRoutes(r *gin.RouterGroup) {

	r.POST("/user", insertAbstractUserHandler)
	r.DELETE("/user", deleteUserHandler)
	r.PATCH("/user", updateAbstractUserHandler)
	// r.GET("/user/:ID", getEmailIDByEmailHandler) // This should not have public access.

	r.POST("/post", insertPostHandler)
	r.GET("/post", getPostsHandler)
	r.PATCH("/post", updatePostHandler)
	r.DELETE("/post", deletePostHandler)
	r.POST("/post/like/:ID", likePostHandler)

	r.POST("/login", loginUserHandler)

	r.POST("/verifyToken", verifyEmail)
	r.POST("/addTTL", insertRandomTokenHandler)

	r.GET("/postTitles", getPostTitlesHandler)

}
