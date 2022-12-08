package views

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/errors"
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AddPostRoutes(r *gin.RouterGroup) {
	r.POST("/post", insertPostHandler)
	r.GET("/post", getPostsHandler)
	r.PATCH("/post", updatePostHandler)
	r.DELETE("/post", deletePostHandler)
	r.POST("/post/like/:ID", likePostHandler)

	r.GET("/postTitles", getPostTitlesHandler)
}

func insertPostHandler(ctx *gin.Context) {
	var post models.PostToCreate

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.PostCreateError)
		return
	}

	if err = controllers.InsertPost(&post, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func getPostsHandler(ctx *gin.Context) {
	var users []models.PostToGet
	title := ctx.Query("title")
	title = strings.ToLower(title)
	if strings.TrimSpace(title) == "" {
		ctx.JSON(http.StatusBadRequest, errors.InvalidTitle)
		return
	}

	if err := controllers.FindPostByPostTitle(&users, &title); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func getPostTitlesHandler(ctx *gin.Context) {
	g := utils.SetLimitFields(ctx)

	var titles []string
	title := ctx.Query("title")
	title = strings.ToLower(title)
	if strings.TrimSpace(title) == "" {
		ctx.JSON(http.StatusBadRequest, errors.InvalidTitle)
		return
	}

	if err := controllers.FindPostTitles(&titles, &title, g.Limit); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, titles)
}

// PATCH, "/post"
func updatePostHandler(ctx *gin.Context) {
	var post models.PostToUpdate
	email, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBind(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.UpdatePostError)
		return
	}
	if err = controllers.UpdatePost(&post, &email); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, utils.SUCCESSFUL)
}

func deletePostHandler(ctx *gin.Context) {
	var post models.PostToUpdate

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBind(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if status, err := controllers.DeletePost(&emailID, &post.Id); err != nil {
		ctx.JSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}

func likePostHandler(ctx *gin.Context) {
	var PostID utils.GeneralID

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBindUri(&PostID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.LikePost(&emailID, &PostID.ID); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}
