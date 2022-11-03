package main

import "github.com/gin-gonic/gin"

func addAPIRoutes(r *gin.RouterGroup) {
	r.GET("/id", GetUsersById)
	r.POST("/id", insertAbstractUserHandler)

	r.POST("/post", insertPostHandler)
	r.GET("/post", getPostsHandler)

	r.PATCH("/post", updatePostHandler)
}
