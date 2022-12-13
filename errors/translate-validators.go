package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Camp  string
	Motiv string
}

func TranslateValidators(e error) gin.H {
	var ve validator.ValidationErrors
	if errors.As(e, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe)}
		}
		return gin.H{"erori": out}

	} else {
		return nil
	}

}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return Required
	case "email":
		return InvalidEmailFormat
	case "pw":
		return PasswordFormat
	case "spacetrim":
		return Required
	case "eqfield":
		return IdenticalFields
	}

	return fe.Error() // default error
}
