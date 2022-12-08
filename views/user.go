package views

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AddUserRoutes(r *gin.RouterGroup) {

	r.POST("/user", insertAbstractUserHandler)
	r.DELETE("/user", deleteUserHandler)
	r.PATCH("/user", updateAbstractUserHandler)

	r.POST("/login", loginUserHandler)

	r.POST("/verifyToken", verifyEmail)

}

func insertAbstractUserHandler(ctx *gin.Context) {
	var absUsr models.AbstractUser

	if err := ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := controllers.EmailExists(absUsr.Email); err != nil {
		ctx.JSON(http.StatusConflict, err.Error())
		return
	}

	verCode, err := utils.SendEmail(&absUsr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: absUsr, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()

	ctx.JSON(http.StatusAccepted, "EmailID sent")

}

func verifyEmail(ctx *gin.Context) {

	var absUsr models.AbstractUser
	if err := ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	verCode := ctx.Query("token")
	utils.MutexVerification.Lock()
	val, ok := utils.CodeToTTL[verCode]
	utils.MutexVerification.Unlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, (&utils.InvalidFieldsError{AffectedField: "token",
			Reason: "Invalid verification code", Location: "query params"}).Error())
		return
	} else if val.AbsUsr.Email == absUsr.Email {

		if err := controllers.InsertAbstractUser(&absUsr); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else if val.Expired() {
		ctx.JSON(http.StatusGone, (&utils.InvalidFieldsError{AffectedField: "token", Reason: "Expired", Location: "query params"}).Error())
		return
	} else {
		ctx.JSON(http.StatusUnauthorized, (&utils.InvalidFieldsError{AffectedField: "token",
			Reason: "Could not match " + "verification code with the provided email", Location: "query params"}).Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}

func updateAbstractUserHandler(ctx *gin.Context) {
	var u models.AbstractUserToUpdate
	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}
	if err = ctx.ShouldBind(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.UpdateAbstractUser(&u, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, utils.SUCCESSFUL)
}

func loginUserHandler(ctx *gin.Context) {
	var u utils.UserLoginFromHeader
	if err := ctx.BindHeader(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	uc, err := utils.DecodeAuth(u.Auth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err = utils.CheckCredentials(&uc); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	var emailID int
	if err = utils.GetEmailIDByEmail(&uc.Email, &emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	sessID := utils.RandomString(utils.SESSION_ID_LENGTH)
	utils.MutexSession.Lock()
	for {
		if _, ok := utils.SessionToEmailID[sessID]; !ok {
			utils.SessionToEmailID[sessID] = utils.CachedLoginSessions{EmailID: emailID,
				SessTTL: time.Now().Add(utils.SESSION_TTL_N_SECONDS)}
			break
		}
		sessID = utils.RandomString(utils.SESSION_ID_LENGTH)
	}
	utils.MutexSession.Unlock()

	if err = utils.CreateSession(&emailID, &sessID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c := http.Cookie{Name: utils.SESSION_COOKIE_NAME, Value: sessID, MaxAge: utils.SESSION_TTL_N_SECONDS, Secure: true}
	ctx.SetSameSite(http.SameSiteNoneMode)
	utils.SetCookieByHTTPCookie(ctx, &c)

	ctx.JSON(http.StatusOK, utils.SUCCESSFUL)
}

// deletes the user from DB if credentials are matched
func deleteUserHandler(ctx *gin.Context) {
	var absUsr models.AbstractUser
	if err := ctx.ShouldBind(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var uc = models.UserCredentials{Email: absUsr.Email, Pass: absUsr.Password}
	if err := utils.CheckCredentials(&uc); err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	if err := controllers.DeleteUser(uc.Email); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}
