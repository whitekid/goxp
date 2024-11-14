package validate

import (
	"github.com/go-playground/validator/v10"

	"github.com/whitekid/goxp/errors"
)

var validate = validator.New()

// Validate shortcuts
func Struct(s interface{}) error              { return validate.Struct(s) }
func Var(field interface{}, tag string) error { return validate.Var(field, tag) }

func Vars(vars ...interface{}) error {
	for i := 0; i < len(vars); i += 2 {
		if err := Var(vars[i], (vars[i+1]).(string)); err != nil {
			return err
		}
	}
	return nil
}

func IsValidationError(err error) bool {
	var verr validator.ValidationErrors

	return errors.As(err, &verr)
}
