package utils

import (
	"bytes"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"example/web-service-gin/errors"
	"example/web-service-gin/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type InvalidFieldsError struct {
	AffectedField string
	Reason        string
	Location      string
}

func (m *InvalidFieldsError) Error() string {
	return fmt.Sprintf("Cannot process <%s> field: <%s>. Reason: <%s>", m.Location, m.AffectedField, m.Reason)
}

func SHA512(text string) string {
	h := sha512.New512_256()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func SetLimitFields(ctx *gin.Context) GeneralQueryFields {

	var g GeneralQueryFields
	_ = ctx.ShouldBind(&g)
	g.SetDefault()
	return g
}

func (g *GeneralQueryFields) SetDefault() {
	if g.Limit == 0 {
		g.Limit = 2000
	}
}

func SetCookieByHTTPCookie(ctx *gin.Context, ck *http.Cookie) {
	ctx.SetCookie(ck.Name, ck.Value, ck.MaxAge, ck.Path, ck.Domain, ck.Secure, ck.HttpOnly)
}

func PopulateConfig(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, &Cred)
	if err != nil {
		panic(err)
	}
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func (ver *VerificationTTL) Expired() bool {
	if ver.TTL.Unix() < time.Now().Unix() {
		return true
	}
	return false

}

func DecodeAuth(auth string) (models.UserCredentials, error) {
	if strings.HasPrefix(auth, "Basic ") {
		sDec, err := b64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return models.UserCredentials{}, err
		}
		name, pass, found := bytes.Cut(sDec, []byte{58}) // Separate by ":"
		if found {
			return models.UserCredentials{Email: string(name), Password: string(pass)}, nil
		}
		return models.UserCredentials{}, &InvalidFieldsError{"Authorization", "Invalid format. Missing colon ", "Basic auth"}
	}
	return models.UserCredentials{}, &InvalidFieldsError{"Authorization", "Can only process Basic Authentication", "Basic auth"}
}

func VerifyWithCookie(ctx *gin.Context) (int, error) {
	cookieVal, err := ctx.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.CookieUnfound)
		return -1, err
	}

	MutexSession.RLock()
	CLS, ok := SessionToEmailID[cookieVal]
	MutexSession.RUnlock()

	if !ok {
		ctx.JSON(http.StatusUnauthorized, errors.SessionExpired)
		return -1, errors.New("-1")
	}

	return CLS.EmailID, nil
}

func SendEmail(email *string, message *string) error {
	_from := Cred.AnonymousGMailName
	pw := Cred.AnonymousGmailPass

	to := []string{*email}

	_host := "smtp.gmail.com"
	p := "587"
	address := _host + ":" + p

	_message := []byte(*message)

	auth := smtp.PlainAuth("", _from, pw, _host)

	if err := smtp.SendMail(address, auth, _from, to, _message); err != nil {
		return err
	}
	return nil
}
