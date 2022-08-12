package customValidators

import (
	"github.com/go-playground/validator/v10"
)

var CustomValidate *validator.Validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ConstructValidationError(err error) []*ErrorResponse {
	errors := make([]*ErrorResponse, 0)

	for _, val := range err.(validator.ValidationErrors) {

		var element ErrorResponse
		element.FailedField = val.StructNamespace()
		element.Tag = val.Tag()
		element.Value = val.Param()

		errors = append(errors, &element)
	}
	return errors

}
