package validate

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

func Struct(req interface{}) (string, error) {
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(req)
	if err != nil {
		if errors.As(err, &validator.ValidationErrors{}) {
			validateErr := err.(validator.ValidationErrors)
			return validateErr.Error(), err
		}
		return "invalid request", err
	}
	return "validation completed", nil
}
