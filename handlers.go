package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

/*
From REQUEST BODY gets the user credentials (email and pass are mandatory). Checks if email exists.
Sends a mail with the verification code. Generates a 10-min code available for mail confirmation
*/
func insertAbstractUserHandler(ctx *gin.Context) {
	var absUsr AbstractUser

	if err = ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = emailExists(absUsr.Email); err != nil {
		ctx.JSON(http.StatusConflict, err.Error())
		return
	}

	verCode, err := SendEmail(&absUsr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ttl := VerificationTTL{absUsr, time.Now().Add(time.Second * VERIFICATION_TTL_N_SECONDS)}
	mutex.Lock()
	codeToTTL[verCode] = ttl
	mutex.Unlock()

	ctx.JSON(http.StatusAccepted, "EmailID sent")

}

func verifyEmail(ctx *gin.Context) {

	var absUsr AbstractUser
	if err = ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	verCode := ctx.Query("token")
	mutex.Lock()
	val, ok := codeToTTL[verCode]
	mutex.Unlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, (&InvalidFieldsError{affectedField: "token", reason: "Invalid verification code", location: "query params"}).Error())
		return
	} else if val.AbsUsr.Email == absUsr.Email {

		if err = insertAbstractUser(&absUsr); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else if val.expired() {
		ctx.JSON(http.StatusGone, (&InvalidFieldsError{affectedField: "token", reason: "Expired", location: "query params"}).Error())
		return
	} else {
		ctx.JSON(http.StatusUnauthorized, (&InvalidFieldsError{affectedField: "token", reason: "Could not match " +
			"verification code with the provided email", location: "query params"}).Error())
		return
	}
	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
}

// POST, "/post"
func insertPostHandler(ctx *gin.Context) {
	var post PostToCreate

	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = insertPost(&post, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, SUCCESSFUL)

}

func getPostsHandler(ctx *gin.Context) {
	var users []PostToGet
	title := ctx.Query("title")
	title = strings.ToLower(title)

	if err = findPostByPostTitle(&users, &title); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func getPostTitlesHandler(ctx *gin.Context) {
	g := setLimitFields(ctx)
	var titles []string
	title := ctx.Query("title")
	title = strings.ToLower(title)

	if err = findPostTitles(&titles, &title, g.Limit); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, titles)
}

// PATCH, "/post"
func updatePostHandler(ctx *gin.Context) {
	var post PostToUpdate
	email, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBind(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err = updatePost(&post, &email); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, SUCCESSFUL)
}

func updateAbstractUserHandler(ctx *gin.Context) {
	var u AbstractUserToUpdate
	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.ShouldBind(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = updateAbstractUser(&u, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, SUCCESSFUL)
}

func loginUserHandler(ctx *gin.Context) {
	var u UserLoginFromHeader
	if err = ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	uc, err := decodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = checkCredentials(&uc); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	var emailID int
	if err = getEmailIDByEmail(&uc.Email, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	sessID := randomString(SESSION_ID_LENGTH)
	mutex.Lock()
	for {
		if _, ok := sessionToEmailID[sessID]; !ok {
			sessionToEmailID[sessID] = CachedLoginSessions{EmailID: emailID, SessTTL: time.Now().Add(SESSION_TTL_N_SECONDS)}
			break
		}
		sessID = randomString(SESSION_ID_LENGTH)
	}
	mutex.Unlock()

	if err = createSession(&emailID, &sessID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c := http.Cookie{Name: SESSION_COOKIE_NAME, Value: sessID, MaxAge: SESSION_TTL_N_SECONDS}
	setCookieByHTTPCookie(ctx, &c)

	ctx.JSON(http.StatusOK, SUCCESSFUL)
}

// Only used for testing purposes (Generates a random token without sending email)
func insertRandomTokenHandler(ctx *gin.Context) {
	var absUsr AbstractUser

	if err = ctx.ShouldBind(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	verCode := randomString(VERIFICATION_CODE_LENGTH)
	ttl := VerificationTTL{absUsr, time.Now().Add(VERIFICATION_TTL_N_SECONDS * time.Second)}

	mutex.Lock()
	codeToTTL[verCode] = ttl
	mutex.Unlock()

	ctx.JSON(http.StatusAccepted, verCode)

}

// deletes the user from DB if credentials are matched
func deleteUserHandler(ctx *gin.Context) {
	var absUsr AbstractUser
	if err = ctx.ShouldBind(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var uc = UserCredentials{absUsr.Email, absUsr.Password}
	if err = checkCredentials(&uc); err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	if err = deleteUser(uc.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
}

func deletePostHandler(ctx *gin.Context) {
	var post PostToUpdate

	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBind(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if status, err := deletePost(&emailID, &post.Id); err != nil {
		ctx.JSON(status, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
}

func likePostHandler(ctx *gin.Context) {
	var PostID GeneralID

	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBindUri(&PostID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = likePost(&emailID, &PostID.ID); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
}

// This should never be publicly accessible.
//func getEmailIDByEmailHandler(ctx *gin.Context) {
//	var email GeneralString
//	if err = ctx.ShouldBindUri(&email); err != nil {
//		ctx.JSON(http.StatusBadRequest, err.Error())
//		return
//	}
//
//	var emailID int
//	if err = getEmailIDByEmail(&email.Value, &emailID); err != nil {
//		ctx.JSON(http.StatusBadRequest, err.Error())
//		return
//	}
//	ctx.JSON(http.StatusOK, emailID)
//}

func insertCommentHandler(ctx *gin.Context) {
	var comm CommentToCreate

	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.BindJSON(&comm); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = insertComment(&comm, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, SUCCESSFUL)

}

func likeCommentHandler(ctx *gin.Context) {
	var commentID GeneralID

	emailID, err := verifyWithCookie(ctx)
	if err != nil {
		return
	}

	if err = ctx.ShouldBindUri(&commentID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = likeComment(&emailID, &commentID.ID); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, SUCCESSFUL)
}

func updateCommentHandler(ctx *gin.Context) {
	var comm CommentToCreate
	var commentID GeneralID

	emailID, err := verifyWithCookie(ctx)
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

	if err = updateComment(&comm, &emailID, &commentID.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, SUCCESSFUL)

}

func getCommentHandler(ctx *gin.Context) {
	var postID GeneralID
	var comments []CommentToGet
	if err = ctx.ShouldBindUri(&postID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = getCommentsFromPost(&postID.ID, &comments); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, comments)
}
