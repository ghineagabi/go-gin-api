package views

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/errors"
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddCommentRoutes(r *gin.RouterGroup) {
	r.GET("/post/:ID/comment", getCommentHandler)
	r.POST("/post/comment", insertCommentHandler)
	r.POST("/post/comment/like/:ID", likeCommentHandler)
	r.POST("/post/comment/respond/:ID", insertRespondToCommentHandler)
	r.PATCH("/post/comment/:ID", updateCommentHandler)
}

func insertCommentHandler(ctx *gin.Context) {
	var comm models.CommentToCreate

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&comm); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.CommentError)
		return
	}

	if err = controllers.InsertComment(&comm, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.InsertCommentError)
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func insertRespondToCommentHandler(ctx *gin.Context) {
	var comm models.CommentToCreate
	var commID utils.GeneralID

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&comm); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.CommentError)
		return
	}
	if err = ctx.ShouldBindUri(&commID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.InsertRespondToComment(&comm, &emailID, &commID.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func likeCommentHandler(ctx *gin.Context) {
	var commentID utils.GeneralID

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBindUri(&commentID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.LikeComment(&emailID, &commentID.ID); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}

func updateCommentHandler(ctx *gin.Context) {
	var comm models.CommentToCreate
	var commentID utils.GeneralID

	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&comm); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err = ctx.ShouldBindUri(&commentID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.UpdateComment(&comm, &emailID, &commentID.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func getCommentHandler(ctx *gin.Context) {
	var postID utils.GeneralID
	var comments []models.CommentToGet

	id := ctx.Query("commentID")
	commentID, err := strconv.Atoi(id)
	if err != nil {
		commentID = 0
	}

	if err = ctx.ShouldBindUri(&postID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.GetCommentsFromPost(&postID.ID, &comments, &commentID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, comments)
}
