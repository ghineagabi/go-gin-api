package main

import (
	"bytes"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
	"unicode"
)

func onlyUnicode(s string) bool {
	if s = strings.TrimSpace(s); s == "" {
		return false
	}
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func SendEmail(absUsr *AbstractUser) (string, error) {
	_from := cred.AnonymousGMailName
	pw := cred.AnonymousGmailPass

	to := []string{absUsr.Email}

	_host := "smtp.gmail.com"
	p := "587"
	address := _host + ":" + p
	verCode := randomString(VERIFICATION_CODE_LENGTH)
	subject := "Subject: This is the subject of the mail\r\n" + "\r\n"
	body := "Hello. We see you are trying to create an account. In order to continue, you need to validate your" +
		" email address. Please use this verification code to continue: " + verCode + "\r\n"

	message := []byte(subject + body)

	auth := smtp.PlainAuth("", _from, pw, _host)

	err = smtp.SendMail(address, auth, _from, to, message)
	if err != nil {
		return "", err
	}
	return verCode, nil
}

func (g *GeneralQueryFields) SetDefault() {
	if g.Limit == 0 {
		g.Limit = 2000
	}
}

func SHA512(text string) string {
	h := sha512.New512_256()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func decodeAuth(auth string) (UserCredentials, error) {
	if strings.HasPrefix(auth, "Basic ") {
		sDec, err := b64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return UserCredentials{}, err
		}
		name, pass, found := bytes.Cut(sDec, []byte{58}) // Separate by ":"
		if found {
			return UserCredentials{string(name), string(pass)}, nil
		}
	}
	return UserCredentials{}, &InvalidFieldsError{"Authorization", "Can only process Basic Authentication", "Basic auth"}
}

func setCookieByHTTPCookie(ctx *gin.Context, ck *http.Cookie) {
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

	err = json.Unmarshal(byteValue, &cred)
	if err != nil {
		panic(err)
	}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func (ver *VerificationTTL) expired() bool {
	if ver.TTL.Unix() > time.Now().Unix() {
		return true
	}
	return false

}

func verifyWithCookie(ctx *gin.Context) (int, error) {
	cookieVal, err := ctx.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return -1, err
	}

	mutex.Lock()
	CLS, ok := sessionToEmailID[cookieVal]
	mutex.Unlock()

	if !ok {
		err = &InvalidFieldsError{affectedField: "Cookie", reason: "Expired Session", location: "Headers"}
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return -1, err
	}

	return CLS.EmailID, nil
}

func setLimitFields(ctx *gin.Context) GeneralQueryFields {
	var g GeneralQueryFields
	_ = ctx.ShouldBind(&g)
	g.SetDefault()
	return g
}
