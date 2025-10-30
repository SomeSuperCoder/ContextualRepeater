package validators

import (
	"github.com/go-playground/validator/v10"
)

type AccessValidator struct {
	validator *validator.Validate
}

func NewDefaultValidator() *AccessValidator {
	v := validator.New()
	av := &AccessValidator{
		validator: v,
	}

	return av
}

func (av *AccessValidator) ValidateRequest(r any) error {
	return av.validator.Struct(r)
}
