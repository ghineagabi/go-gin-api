package main

import (
	"bytes"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"unicode"
)

func onlyUnicodeWithoutSpaces(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII || c == 32 {
			return false
		}
	}
	return true
}

func stringToStringArray(str string, field string) ([]string, error) {
	if len(str) == 0 {
		x := make([]string, 0)
		return x, &InvalidFieldsError{affectedField: field}
	}
	str = strings.ReplaceAll(str, " ", "")
	return strings.Split(str, ","), nil
}

func SendEmail(absUsr *AbstractUser) error {
	_from := cred.AnonymousGMailName
	pw := cred.AnonymousGmailPass

	to := []string{absUsr.Email}

	_host := "smtp.gmail.com"
	p := "587"
	address := _host + ":" + p

	subject := "Subject: This is the subject of the mail\n"
	body := "This is the body of the mail"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", _from, pw, _host)

	err = smtp.SendMail(address, auth, _from, to, message)
	if err != nil {
		return err
	}
	return nil
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
		name, pass, found := bytes.Cut(sDec, []byte{58})
		if found {
			return UserCredentials{string(name), string(pass)}, nil
		}
	}
	return UserCredentials{}, &InvalidFieldsError{"Authorization", "Can only process Basic Authentication"}
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
