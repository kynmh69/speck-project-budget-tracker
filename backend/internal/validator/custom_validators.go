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
