package main

import (
	"example/web-service-gin/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var nonEmpty validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	} else {
		return utils.OnlyUnicode(v)
	}
}

var validPass validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	} else {
		return utils.ValidatePassword(v)
	}
}

func addValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		if err := v.RegisterValidation("spacetrim", nonEmpty); err != nil {
			return
		}
		if err := v.RegisterValidation("pw", validPass); err != nil {
			return
		}
	}
}
