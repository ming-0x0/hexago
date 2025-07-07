package customer

import (
	"github.com/ming-0x0/hexago/internal/modules/customer/domain/service_type"
	"github.com/ming-0x0/hexago/internal/shared/domain/email"
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
