package main

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var nonEmpty validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	} else {
		return onlyUnicode(v)
	}
}

func addValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("spacetrim", nonEmpty)
		if err != nil {
			return
		}
		v.RegisterStructValidation(UserStructLevelValidation, AbstractUser{})
	}
}

func UserStructLevelValidation(sl validator.StructLevel) {
	absUsr := sl.Current().Interface().(AbstractUser)

	if len(absUsr.Password) < 8 || len(absUsr.Password) > 50 {
		sl.ReportError(absUsr.Password, "password", "password", "pass-len", "")
	}

}
