package status

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/hexago/internal/shared/errors"
)

var (
	Replied   = Status{value: 1}
	Unreplied = Status{value: 2}
)

//go:generate sh -c "$(go list -m -f '{{.Dir}}')/bin/accessor -type=Status"
type Status struct {
	value int64
}

func New(value int64) (*Status, error) {
	s := &Status{value: value}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Status) validate() error {
	err := validation.ValidateStruct(s,
		validation.Field(
			&s.value,
			validation.Required,
			validation.In(Replied.value, Unreplied.value),
		),
	)
	if err != nil {
		return errors.NewDomainError(errors.Validation, err.Error())
	}

	return nil
}
