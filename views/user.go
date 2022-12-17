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
	r.POST("/verifyEmail", verifyEmail)
	r.DELETE("/user", deleteUserHandler)
	r.PATCH("/user", updateAbstractUserHandler)

	r.POST("/login", loginUserHandler)
	r.POST("/logout", logoutUserHandler)

	r.POST("/resetPassword", resetPasswordHandler)
	r.POST("/verifyForgotPassword", verifyForgotPasswordHandler)
	r.POST("/forgotPassword", forgotPasswordHandler)

}

func insertAbstractUserHandler(ctx *gin.Context) {
	var absUsr models.AbstractUser

	if err := ctx.BindJSON(&absUsr); err != nil {
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
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
		ctx.JSON(http.StatusBadRequest, errors.EmailSendingError)
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: absUsr, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()

	ctx.JSON(http.StatusAccepted, "Email trimis.")

}

func verifyEmail(ctx *gin.Context) {

	verCode := ctx.Query("token")
	utils.MutexVerification.RLock()
	val, ok := utils.CodeToTTL[verCode]
	utils.MutexVerification.RUnlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, errors.InvalidToken)
		return
	} else if val.Expired() {
		ctx.JSON(http.StatusGone, errors.ExpiredToken)
		return
	}

	if err := controllers.InsertAbstractUser(&val.AbsUsr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	utils.MutexVerification.Lock()
	delete(utils.CodeToTTL, verCode)
	utils.MutexVerification.Unlock()

	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)
	return

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
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
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

	if err := ctx.BindJSON(&uc); err != nil {
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
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

func resetPasswordHandler(ctx *gin.Context) {
	emailID, err := utils.VerifyWithCookie(ctx)
	if err != nil {
		return
	}

	var RP models.ResetPassword
	if err = ctx.BindJSON(&RP); err != nil {
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = controllers.UpdatePassword(&emailID, &RP.ConfirmPassword, &RP.OldPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func forgotPasswordHandler(ctx *gin.Context) {
	var GE utils.GeneralEmail
	var emailID int

	if err := ctx.BindJSON(&GE); err != nil {
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	email := &GE.Email

	verCode := utils.RandomString(utils.VERIFICATION_WHEN_FORGOT_LENGTH)
	subject := "Subject: Beautyfinder password reset\r\n" + "\r\n"
	body := "Hello. We see you are trying to reset your password. In order to continue, you need to validate your" +
		" email address. Please use this verification code to continue: " + verCode + "\r\n"
	message := subject + body

	if err := utils.GetEmailIDByEmail(email, &emailID); err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err := utils.SendEmail(email, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ttl := utils.VerificationTTL{AbsUsr: models.AbstractUser{Email: *email}, TTL: time.Now().Add(time.Second * utils.VERIFICATION_TTL_N_SECONDS)}
	utils.MutexVerification.Lock()
	utils.CodeToTTL[verCode] = ttl
	utils.MutexVerification.Unlock()
}

func verifyForgotPasswordHandler(ctx *gin.Context) {

	var FP models.ForgotPassword
	if err := ctx.BindJSON(&FP); err != nil {
		if valErr := errors.TranslateValidators(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	newPass := &FP.ConfirmPassword

	verCode := ctx.Query("token")

	utils.MutexVerification.RLock()
	val, ok := utils.CodeToTTL[verCode]
	utils.MutexVerification.RUnlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, errors.InvalidToken)
		return
	} else if val.Expired() {
		ctx.JSON(http.StatusGone, errors.ExpiredToken)
		return
	}

	if err := controllers.UpdatePasswordByEmail(&val.AbsUsr.Email, newPass); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}

func logoutUserHandler(ctx *gin.Context) {
	cookieVal, err := ctx.Cookie(utils.SESSION_COOKIE_NAME)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.CookieUnfound)
		return
	}

	if _, ok := utils.SessionToEmailID[cookieVal]; !ok {
		ctx.JSON(http.StatusGone, errors.CookieValueUnfound)
		return
	}

	delete(utils.SessionToEmailID, cookieVal)
	if err = utils.DeleteSession(&cookieVal); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusAccepted, utils.SUCCESSFUL)

}
