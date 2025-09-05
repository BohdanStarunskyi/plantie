package utils

import (
	"plant-reminder/constants"

	"github.com/go-playground/validator"
)

var Validate = func() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("validrepeattype", validateRepeatType)

	return validate
}()

func validateRepeatType(fl validator.FieldLevel) bool {
	repeatType := fl.Field().Interface().(constants.RepeatType)

	switch repeatType {
	case constants.RepeatDaily, constants.RepeatWeekly, constants.RepeatMonthly:
		return true
	default:
		return false
	}
}
