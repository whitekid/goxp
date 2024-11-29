package validate

import (
	"github.com/go-playground/validator/v10"

	"github.com/whitekid/goxp/errors"
)

var validate = validator.New()

// Validate shortcuts
func Struct(s any) error {
	if err := validate.Struct(s); err != nil {
		return errors.Errorf(err, "validation failed")
	}

	return nil
}

func Var(field any, tag string) error {
	if err := validate.Var(field, tag); err != nil {
		return errors.Errorf(err, "validation failed")
	}

	return nil
}

func Vars(vars ...any) error {
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
