package validation

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func NotBeforeNowValidator(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(*time.Time)
	if !ok {
		return false
	}
	return (*date).Sub(time.Now()) > -24*time.Hour
}
