package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ApiError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func TranslateValidators(e error) gin.H {
	var ve validator.ValidationErrors
	if errors.As(e, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			_field := fe.Field()
			firstLetterToLowerCase := strings.ToLower(string(_field[0]))
			out[i] = ApiError{firstLetterToLowerCase + _field[1:], msgForTag(fe)}
		}
		return gin.H{"errors": out}

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
	case "nefield":
		return DifferentFields
	}

	return fe.Error() // default error
}
