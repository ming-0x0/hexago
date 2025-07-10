package customer

import (
	"github.com/ming-0x0/hexago/domain/service_type"
	"github.com/ming-0x0/hexago/shared/domain/email"
	"github.com/ming-0x0/hexago/shared/undefined"
)

type ID string

//go:generate accessor -type=Customer
type Customer struct {
	id           ID
	customerName string
	email        email.Email
	phoneNumber  string
	companyName  undefined.Undefined[string]
	message      undefined.Undefined[string]
	note         undefined.Undefined[string]
	serviceType  service_type.ServiceType
}
