package main

import (
	"github.com/gin-gonic/gin"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"net/http"
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

func (g *GeneralQueryFields) SetDefault() {
	if g.Limit == 0 {
		g.Limit = 2000
	}
}

var toValidate = map[string]string{
	"aud": "api://default",
	"cid": os.Getenv("OKTA_CLIENT_ID"),
}

func verify(ctx *gin.Context) bool {
	status := true
	token := ctx.Request.Header.Get("Authorization")

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		verifierSetup := jwtverifier.JwtVerifier{
			Issuer:           "https://" + os.Getenv("OKTA_DOMAIN") + "/oauth2/default",
			ClaimsToValidate: toValidate,
		}
		verifier := verifierSetup.New()
		_, err = verifier.VerifyAccessToken(token)
		if err != nil {
			ctx.String(http.StatusForbidden, err.Error())
			print(err.Error())
			status = false
		}
	} else {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		status = false
	}
	return status
}
