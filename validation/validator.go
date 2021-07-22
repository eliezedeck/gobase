package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var V = validator.New()

var (
	// This is based on https://stackoverflow.com/questions/28020805/regex-validate-correct-iso8601-date-string-with-time
	// but only for the date part
	dateISO8601RegexStr = "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)$"
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
