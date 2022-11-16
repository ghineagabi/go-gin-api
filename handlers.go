package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsersById(ctx *gin.Context) {
	ids := ctx.Query("ids")
	arrayIDs, err := stringToStringArray(ids, "ids")
	if err != nil || ids == "" {
	}
	users, err := findUsersByID(arrayIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	ctx.JSON(http.StatusOK, users)

}

func insertAbstractUserHandler(ctx *gin.Context) {
	var absUsr AbstractUser
	err = ctx.BindJSON(&absUsr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = SendEmail(&absUsr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = insertAbstractUser(&absUsr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, "Successful")

}

func insertPostHandler(ctx *gin.Context) {
	var post Post
	err = ctx.BindJSON(&post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = insertPost(&post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
	ctx.JSON(http.StatusAccepted, "Successful")

}

func getPostsHandler(ctx *gin.Context) {
	var g GeneralQueryFields
	var users []PostInfo
	title := ctx.Query("title")
	_ = ctx.ShouldBind(&g)
	g.SetDefault()
	err = findNameByPostTitle(&users, title)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func updatePostHandler(ctx *gin.Context) {
	var u ToUpdatePost
	err = ctx.ShouldBind(&u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = updatePost(&u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Successful")
}

func updateAbstractUserHandler(ctx *gin.Context) {
	var u AbstractUser
	err = ctx.ShouldBind(&u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, "Successful")
}

func loginUserHandler(ctx *gin.Context) {
	var s AbstractUserSession
	var u UserLoginFromHeader
	if err = ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	uc, err := decodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}
	err = checkCredentials(&uc)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	err = createSession(&s, &uc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c := http.Cookie{Name: "sess-id", Value: s.Id, MaxAge: 60 * 60 * 24}
	setCookieByHTTPCookie(ctx, &c)
	ctx.JSON(http.StatusOK, "Successful login and session creation")
}
