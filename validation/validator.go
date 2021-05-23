package validation

import (
	"github.com/go-playground/validator/v10"
)

var V = validator.New()

// ValidatorErrorMessage returns the first error message after validation; if it's not an instance of ValidationErrors
// then it will just return the normal .Error().
func ValidatorErrorMessage(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		first := ve[0]
		return first.Error()
	}
	return err.Error()
}
