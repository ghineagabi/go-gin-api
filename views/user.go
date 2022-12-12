package views

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/errors"
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

	r.POST("/verifyRegistrationToken", verifyEmail)
	r.POST("/resetPassword", resetPasswordWhenLoggedInHandler)
	r.POST("/verifyResetPassword", VerifyResetPasswordHandler)
	r.POST("/user/logout", logoutUserhandler)

}

func insertAbstractUserHandler(ctx *gin.Context) {
	var absUsr models.AbstractUser

	if err := ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if err := controllers.EmailExists(absUsr.Email); err != nil {
		ctx.JSON(http.StatusConflict, errors.EmailAlreadyExists)
		return
	}

	verCode := utils.RandomString(utils.VERIFICATION_CODE_LENGTH)
	subject := "Subject: Inregistrare beautyfinder\r\n" + "\r\n"
	body := "Buna ziua. Ati aplicat pentru crearea unui cont pe beautyfinder.ro. Pentru a putea continua, trebuie sa" +
		" validati aceasta adresa de email. Folositi acest cod pentru a putea continua: " + verCode + "\r\n"
	message := subject + body

	if err := utils.SendEmail(&absUsr.Email, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.MailSendingError)
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: absUsr, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()

	ctx.JSON(http.StatusAccepted, "Email trimis.")

}

func verifyEmail(ctx *gin.Context) {

	var absUsr models.AbstractUser
	if err := ctx.BindJSON(&absUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.InvalidUserFields)
		return
	}

	verCode := ctx.Query("token")
	utils.MutexVerification.Lock()
	val, ok := utils.CodeToTTL[verCode]
	utils.MutexVerification.Unlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, errors.InvalidToken)
		return
	} else if val.Expired() {
		ctx.JSON(http.StatusGone, errors.ExpiredToken)
		return
	} else if val.AbsUsr.Email == absUsr.Email {
		if err := controllers.InsertAbstractUser(&absUsr); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

	} else {
		ctx.JSON(http.StatusUnauthorized, errors.TokenEmailMismatch)
		return
	}

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

	var emailID int
	if err = utils.CheckCredentials(&uc, &emailID); err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.InvalidLoginCredentials)
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
	utils.SetCookieByHTTPCookie(ctx, &c)

	ctx.JSON(http.StatusOK, utils.SUCCESSFUL)
}

// deletes the user from DB if credentials are matched
func deleteUserHandler(ctx *gin.Context) {
	var uc models.UserCredentials
	if err := ctx.ShouldBind(&uc); err != nil {
		ctx.JSON(http.StatusBadRequest, errors.RequiredEmailPass)
		return
	}

	var emailID int
	if err := utils.CheckCredentials(&uc, &emailID); err != nil {
		ctx.JSON(http.StatusForbidden, errors.InvalidLoginCredentials)
		return
	}

	if err := controllers.DeleteUser(&emailID); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
}

func resetPasswordWhenLoggedInHandler(ctx *gin.Context) {
	var email string
	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	verCode := utils.RandomString(utils.VERIFICATION_CODE_LENGTH)
	subject := "Subject: Beautyfinder password reset\r\n" + "\r\n"
	body := "Hello. We see you are trying to reset your password. In order to continue, you need to validate your" +
		" email address. Please use this verification code to continue: " + verCode + "\r\n"
	message := subject + body

	if err = controllers.GetEmailByEmailID(&emailID, &email); err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	if err := utils.SendEmail(&email, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: models.AbstractUser{Email: email}, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()
}

func resetPasswordWhenForgotHandler(ctx *gin.Context) {
	var email utils.GeneralEmail
	err := ctx.BindJSON(&email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.InvalidEmailFormat)
		return
	}

	verCode := utils.RandomString(utils.VERIFICATION_CODE_LENGTH)
	subject := "Subject: Beautyfinder password reset\r\n" + "\r\n"
	body := "Hello. We see you are trying to reset your password. In order to continue, you need to validate your" +
		" email address. Please use this verification code to continue: " + verCode + "\r\n"
	message := subject + body

	if err := utils.SendEmail(&email.Email, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: models.AbstractUser{Email: email.Email}, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()
}

// TODO: get the password from the header, not from the query params. Also, see if can validate on field level
func VerifyResetPasswordHandler(ctx *gin.Context) {
	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	newPass := ctx.Query("password")
	verCode := ctx.Query("token")

	utils.MutexVerification.Lock()
	val, ok := utils.CodeToTTL[verCode]
	utils.MutexVerification.Unlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, (&utils.InvalidFieldsError{AffectedField: "token",
			Reason: "Invalid verification code", Location: "query params"}).Error())
		return
	} else if val.Expired() {
		ctx.JSON(http.StatusGone, (&utils.InvalidFieldsError{AffectedField: "token", Reason: "Expired", Location: "query params"}).Error())
		return
	} else {
		if err = controllers.UpdatePassword(&emailID, &newPass); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
	}

}

func logoutUserhandler(ctx *gin.Context) {
	cookieVal, err := ctx.Cookie(utils.SESSION_COOKIE_NAME)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.CookieUnfound)
		return
	}

	if _, ok := utils.SessionToEmailID[cookieVal]; !ok {
		ctx.JSON(http.StatusGone, errors.CookieUnfound)
		return
	}

	delete(utils.SessionToEmailID, cookieVal)
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}
