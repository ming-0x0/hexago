package customer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ming-0x0/hexago/internal/customer/domain/service_type"
	"github.com/ming-0x0/hexago/internal/customer/domain/status"
	"github.com/ming-0x0/hexago/internal/shared/domain/email"
	"github.com/ming-0x0/hexago/internal/shared/errors"
	"github.com/ming-0x0/hexago/internal/shared/undefined"
	"github.com/oklog/ulid/v2"
)

const (
	maxCustomerNameLength = 255
	maxCompanyNameLength  = 255
	maxPhoneNumberLength  = 10
	maxMessageLength      = 1000
	maxNoteLength         = 1000
)

type ID string

//go:generate sh -c "$(go list -m -f '{{.Dir}}')/bin/accessor -type=Customer"
type Customer struct {
	id           ID
	customerName string
	email        email.Email
	phoneNumber  string
	companyName  undefined.Undefined[string]
	message      undefined.Undefined[string]
	note         undefined.Undefined[string]
	serviceType  service_type.ServiceType
	status       status.Status
}

func New(
	customerName string,
	email email.Email,
	phoneNumber string,
	companyName undefined.Undefined[string],
	message undefined.Undefined[string],
	note undefined.Undefined[string],
	serviceType service_type.ServiceType,
	status status.Status,
) (*Customer, error) {
	customer := &Customer{
		id:           ID(ulid.Make().String()),
		customerName: customerName,
		email:        email,
		phoneNumber:  phoneNumber,
		companyName:  companyName,
		message:      message,
		note:         note,
		serviceType:  serviceType,
		status:       status,
	}

	if err := customer.validate(); err != nil {
		return nil, err
	}

	return customer, nil
}

func FromRepository(
	id ID,
	customerName string,
	email email.Email,
	phoneNumber string,
	companyName undefined.Undefined[string],
	message undefined.Undefined[string],
	note undefined.Undefined[string],
	serviceType service_type.ServiceType,
	status status.Status,
) (*Customer, error) {
	customer := &Customer{
		id:           id,
		customerName: customerName,
		email:        email,
		phoneNumber:  phoneNumber,
		companyName:  companyName,
		message:      message,
		note:         note,
		serviceType:  serviceType,
		status:       status,
	}

	if err := customer.validate(); err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *Customer) validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(
			&c.customerName,
			validation.Required,
			validation.Length(1, maxCustomerNameLength),
		),
		validation.Field(
			&c.phoneNumber,
			validation.Required,
			validation.Length(maxPhoneNumberLength, maxPhoneNumberLength),
		),
		validation.Field(
			&c.companyName,
			validation.Length(0, maxCompanyNameLength),
		),
		validation.Field(
			&c.message,
			validation.Length(0, maxMessageLength),
		),
		validation.Field(
			&c.note,
			validation.Length(0, maxNoteLength),
		),
	)
	if err != nil {
		return errors.NewDomainError(errors.Validation, err.Error())
	}

	return nil
}
