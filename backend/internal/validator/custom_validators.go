package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register custom validators here if needed
}

// Validate validates a struct
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// EchoValidator adapts the package validator to echo's Validator interface,
// so handlers can call c.Validate(&req).
type EchoValidator struct{}

// Validate implements the echo.Validator interface.
func (v *EchoValidator) Validate(i interface{}) error {
	return validate.Struct(i)
}
