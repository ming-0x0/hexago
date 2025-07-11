package service_type

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/hexago/internal/shared/errors"
)

var (
	TuyenDung = ServiceType{value: 1}
	LienHe    = ServiceType{value: 2}
	KhoaHoc   = ServiceType{value: 3}
)

//go:generate sh -c "$(go list -m -f '{{.Dir}}')/bin/accessor -type=ServiceType"
type ServiceType struct {
	value int64
}

func New(value int64) (*ServiceType, error) {
	s := &ServiceType{value: value}
	if err := s.validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *ServiceType) validate() error {
	err := validation.ValidateStruct(s,
		validation.Field(
			&s.value,
			validation.Required,
			validation.In(TuyenDung.value, LienHe.value, KhoaHoc.value),
		),
	)
	if err != nil {
		return errors.NewDomainError(errors.Validation, err.Error())
	}

	return nil
}
