package email

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/ming-0x0/hexago/internal/shared/errors"
)

//go:generate sh -c "$(go list -m -f '{{.Dir}}')/bin/accessor -type=Email"
type Email struct {
	value string
}

func New(value string) (*Email, error) {
	e := &Email{
		value: value,
	}

	if err := e.validate(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Email) validate() error {
	err := validation.ValidateStruct(e,
		validation.Field(&e.value, validation.Required, is.Email),
	)
	if err != nil {
		return errors.NewDomainError(errors.Validation, err.Error())
	}

	return nil
}
