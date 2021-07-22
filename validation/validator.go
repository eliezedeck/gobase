package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var V = validator.New()

var (
	dateISO8601RegexStr = "^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])(?:T|\\s)(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])?(Z)?$"
	dateISO8601Regex    = regexp.MustCompile(dateISO8601RegexStr)
)

// ValidatorErrorMessage returns the first error message after validation; if it's not an instance of ValidationErrors
// then it will just return the normal .Error().
func ValidatorErrorMessage(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		first := ve[0]
		return first.Error()
	}
	return err.Error()
}

func init() {
	if err := V.RegisterValidation("date-iso8601", func(fl validator.FieldLevel) bool {
		return dateISO8601Regex.MatchString(fl.Field().String())
	}); err != nil {
		panic(err)
	}
}
