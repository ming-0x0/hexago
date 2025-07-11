package customer

import (
	"strings"
	"testing"

	"github.com/ming-0x0/hexago/internal/customer/domain/service_type"
	"github.com/ming-0x0/hexago/internal/customer/domain/status"
	"github.com/ming-0x0/hexago/internal/shared/domain/email"
	"github.com/ming-0x0/hexago/internal/shared/undefined"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	validEmail, _ := email.New("test@example.com")
	validServiceType, _ := service_type.New(1)
	validStatus, _ := status.New(1)
	undefinedStr := undefined.Undefined[string]{}

	type args struct {
		customerName string
		email        email.Email
		phoneNumber  string
		companyName  undefined.Undefined[string]
		message      undefined.Undefined[string]
		note         undefined.Undefined[string]
		serviceType  service_type.ServiceType
		status       status.Status
	}

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Valid",
			args: args{
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.NoError,
		},
		{
			name: "Invalid_CustomerName",
			args: args{
				customerName: strings.Repeat("a", maxCustomerNameLength+1),
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_CustomerName",
			args: args{
				customerName: strings.Repeat("a", maxCustomerNameLength+1),
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_PhoneNumber",
			args: args{
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  strings.Repeat("1", maxPhoneNumberLength+1),
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_CompanyName",
			args: args{
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefined.New(strings.Repeat("a", maxCompanyNameLength+1)),
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_Message",
			args: args{
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefined.New(strings.Repeat("a", maxMessageLength+1)),
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_Note",
			args: args{
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefined.New(strings.Repeat("a", maxNoteLength+1)),
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(
				tc.args.customerName,
				tc.args.email,
				tc.args.phoneNumber,
				tc.args.companyName,
				tc.args.message,
				tc.args.note,
				tc.args.serviceType,
				tc.args.status,
			)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.args.customerName, got.customerName)
				assert.Equal(t, tc.args.email, got.email)
				assert.Equal(t, tc.args.phoneNumber, got.phoneNumber)
				assert.Equal(t, tc.args.companyName, got.companyName)
				assert.Equal(t, tc.args.message, got.message)
				assert.Equal(t, tc.args.note, got.note)
				assert.Equal(t, tc.args.serviceType, got.serviceType)
				assert.Equal(t, tc.args.status, got.status)
			}
		})
	}
}

func TestFromRepository(t *testing.T) {
	t.Parallel()

	id := ID("customer_id")
	validEmail, _ := email.New("test@example.com")
	validServiceType, _ := service_type.New(1)
	validStatus, _ := status.New(1)
	undefinedStr := undefined.Undefined[string]{}

	type args struct {
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

	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Valid",
			args: args{
				id:           id,
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.NoError,
		},
		{
			name: "Invalid_CustomerName",
			args: args{
				id:           id,
				customerName: strings.Repeat("a", maxCustomerNameLength+1),
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_CustomerName",
			args: args{
				id:           id,
				customerName: strings.Repeat("a", maxCustomerNameLength+1),
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_CustomerName",
			args: args{
				customerName: strings.Repeat("a", maxCustomerNameLength+1),
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_PhoneNumber",
			args: args{
				id:           id,
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  strings.Repeat("1", maxPhoneNumberLength+1),
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_CompanyName",
			args: args{
				id:           id,
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefined.New(strings.Repeat("a", maxCompanyNameLength+1)),
				message:      undefinedStr,
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_Message",
			args: args{
				id:           id,
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefined.New(strings.Repeat("a", maxMessageLength+1)),
				note:         undefinedStr,
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
		{
			name: "Invalid_Note",
			args: args{
				id:           id,
				customerName: "test",
				email:        *validEmail,
				phoneNumber:  "1234567890",
				companyName:  undefinedStr,
				message:      undefinedStr,
				note:         undefined.New(strings.Repeat("a", maxNoteLength+1)),
				serviceType:  *validServiceType,
				status:       *validStatus,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := FromRepository(
				tc.args.id,
				tc.args.customerName,
				tc.args.email,
				tc.args.phoneNumber,
				tc.args.companyName,
				tc.args.message,
				tc.args.note,
				tc.args.serviceType,
				tc.args.status,
			)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.args.id, got.id)
				assert.Equal(t, tc.args.customerName, got.customerName)
				assert.Equal(t, tc.args.email, got.email)
				assert.Equal(t, tc.args.phoneNumber, got.phoneNumber)
				assert.Equal(t, tc.args.companyName, got.companyName)
				assert.Equal(t, tc.args.message, got.message)
				assert.Equal(t, tc.args.note, got.note)
				assert.Equal(t, tc.args.serviceType, got.serviceType)
				assert.Equal(t, tc.args.status, got.status)
			}
		})
	}
}
