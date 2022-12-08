package main

import (
	"example/web-service-gin/views"
	"github.com/gin-gonic/gin"
)

func addAPIRoutes(r *gin.RouterGroup) {
	views.AddPostRoutes(r)
	views.AddUserRoutes(r)
	views.AddCommentRoutes(r)
}
