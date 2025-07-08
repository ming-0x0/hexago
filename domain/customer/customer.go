package customer

import (
	"github.com/ming-0x0/hexago/domain/service_type"
	"github.com/ming-0x0/hexago/shared/domain/email"
)

type ID string

type Customer struct {
	id           ID
	customerName string
	email        email.Email
	phoneNumber  string
	companyName  *string
	message      *string
	note         *string
	serviceType  service_type.ServiceType
}
