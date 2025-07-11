package customer

import (
	"testing"

	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/customer/domain/service_type"
	"github.com/ming-0x0/hexago/internal/customer/domain/status"
	"github.com/ming-0x0/hexago/internal/shared/domain/email"
	"github.com/ming-0x0/hexago/internal/shared/undefined"
	"github.com/stretchr/testify/assert"
)

func TestCustomerRepositoryAdapter_ToDomain(t *testing.T) {
	t.Parallel()

	type args struct {
		e *entity.Customer
	}

	tests := []struct {
		name      string
		args      args
		want      *customer.Customer
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				e: &entity.Customer{
					ID:           "customer_id",
					CustomerName: "test",
					Email:        "test@example.com",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  1,
					Status:       1,
				},
			},
			want: func() *customer.Customer {
				domain, _ := customer.FromRepository(
					"customer_id",
					"test",
					func() email.Email {
						email, _ := email.New("test@example.com")
						return *email
					}(),
					"1234567890",
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					func() service_type.ServiceType {
						serviceType, _ := service_type.New(1)
						return *serviceType
					}(),
					func() status.Status {
						status, _ := status.New(1)
						return *status
					}(),
				)

				return domain
			}(),
			assertion: assert.NoError,
		},
		{
			name: "Invalid_Email",
			args: args{
				e: &entity.Customer{
					ID:           "customer_id",
					CustomerName: "test",
					Email:        "invalid_email",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  1,
					Status:       1,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "Invalid_ServiceType",
			args: args{
				e: &entity.Customer{
					ID:           "customer_id",
					CustomerName: "test",
					Email:        "test@example.com",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  4,
					Status:       1,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "Invalid_Status",
			args: args{
				e: &entity.Customer{
					ID:           "customer_id",
					CustomerName: "test",
					Email:        "test@example.com",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  1,
					Status:       3,
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := &CustomerRepositoryAdapter{}
			got, err := a.ToDomain(tc.args.e)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestCustomerRepositoryAdapter_ToDomains(t *testing.T) {
	t.Parallel()

	type args struct {
		es []*entity.Customer
	}

	tests := []struct {
		name      string
		args      args
		want      []*customer.Customer
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				es: []*entity.Customer{
					{
						ID:           "customer_id_1",
						CustomerName: "test1",
						Email:        "test1@example.com",
						PhoneNumber:  "1234567890",
						CompanyName:  undefined.Undefined[string]{},
						Message:      undefined.Undefined[string]{},
						Note:         undefined.Undefined[string]{},
						ServiceType:  1,
						Status:       1,
					},
					{
						ID:           "customer_id_2",
						CustomerName: "test2",
						Email:        "test2@example.com",
						PhoneNumber:  "0987654321",
						CompanyName:  undefined.Undefined[string]{},
						Message:      undefined.Undefined[string]{},
						Note:         undefined.Undefined[string]{},
						ServiceType:  2,
						Status:       2,
					},
				},
			},
			want: func() []*customer.Customer {
				domain1, _ := customer.FromRepository(
					"customer_id_1",
					"test1",
					func() email.Email {
						email, _ := email.New("test1@example.com")
						return *email
					}(),
					"1234567890",
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					func() service_type.ServiceType {
						serviceType, _ := service_type.New(1)
						return *serviceType
					}(),
					func() status.Status {
						status, _ := status.New(1)
						return *status
					}(),
				)
				domain2, _ := customer.FromRepository(
					"customer_id_2",
					"test2",
					func() email.Email {
						email, _ := email.New("test2@example.com")
						return *email
					}(),
					"0987654321",
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					undefined.Undefined[string]{},
					func() service_type.ServiceType {
						serviceType, _ := service_type.New(2)
						return *serviceType
					}(),
					func() status.Status {
						status, _ := status.New(2)
						return *status
					}(),
				)
				return []*customer.Customer{domain1, domain2}
			}(),
			assertion: assert.NoError,
		},
		{
			name: "Invalid_Email_In_One_Entity",
			args: args{
				es: []*entity.Customer{
					{
						ID:           "customer_id_1",
						CustomerName: "test1",
						Email:        "test1@example.com",
						PhoneNumber:  "1234567890",
						CompanyName:  undefined.Undefined[string]{},
						Message:      undefined.Undefined[string]{},
						Note:         undefined.Undefined[string]{},
						ServiceType:  1,
						Status:       1,
					},
					{
						ID:           "customer_id_2",
						CustomerName: "test2",
						Email:        "invalid_email",
						PhoneNumber:  "0987654321",
						CompanyName:  undefined.Undefined[string]{},
						Message:      undefined.Undefined[string]{},
						Note:         undefined.Undefined[string]{},
						ServiceType:  2,
						Status:       2,
					},
				},
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := &CustomerRepositoryAdapter{}
			got, err := a.ToDomains(tc.args.es)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestCustomerRepositoryAdapter_ToEntities(t *testing.T) {
	t.Parallel()

	type args struct {
		ds []*customer.Customer
	}

	tests := []struct {
		name      string
		args      args
		want      []*entity.Customer
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ds: func() []*customer.Customer {
					domain1, _ := customer.FromRepository(
						"customer_id_1",
						"test1",
						func() email.Email {
							email, _ := email.New("test1@example.com")
							return *email
						}(),
						"1234567890",
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						func() service_type.ServiceType {
							serviceType, _ := service_type.New(1)
							return *serviceType
						}(),
						func() status.Status {
							status, _ := status.New(1)
							return *status
						}(),
					)
					domain2, _ := customer.FromRepository(
						"customer_id_2",
						"test2",
						func() email.Email {
							email, _ := email.New("test2@example.com")
							return *email
						}(),
						"0987654321",
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						func() service_type.ServiceType {
							serviceType, _ := service_type.New(2)
							return *serviceType
						}(),
						func() status.Status {
							status, _ := status.New(2)
							return *status
						}(),
					)
					return []*customer.Customer{domain1, domain2}
				}(),
			},
			want: []*entity.Customer{
				{
					ID:           "customer_id_1",
					CustomerName: "test1",
					Email:        "test1@example.com",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  1,
					Status:       1,
				},
				{
					ID:           "customer_id_2",
					CustomerName: "test2",
					Email:        "test2@example.com",
					PhoneNumber:  "0987654321",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  2,
					Status:       2,
				},
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := &CustomerRepositoryAdapter{}
			got, err := a.ToEntities(tc.args.ds)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
func TestCustomerRepositoryAdapter_ToEntity(t *testing.T) {
	t.Parallel()

	type args struct {
		d *customer.Customer
	}

	tests := []struct {
		name      string
		args      args
		want      *entity.Customer
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				d: func() *customer.Customer {
					domain, _ := customer.FromRepository(
						"customer_id",
						"test",
						func() email.Email {
							email, _ := email.New("test@example.com")
							return *email
						}(),
						"1234567890",
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						undefined.Undefined[string]{},
						func() service_type.ServiceType {
							serviceType, _ := service_type.New(1)
							return *serviceType
						}(),
						func() status.Status {
							status, _ := status.New(1)
							return *status
						}(),
					)

					return domain
				}(),
			},
			want: func() *entity.Customer {
				return &entity.Customer{
					ID:           "customer_id",
					CustomerName: "test",
					Email:        "test@example.com",
					PhoneNumber:  "1234567890",
					CompanyName:  undefined.Undefined[string]{},
					Message:      undefined.Undefined[string]{},
					Note:         undefined.Undefined[string]{},
					ServiceType:  1,
					Status:       1,
				}
			}(),
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := &CustomerRepositoryAdapter{}
			got, err := a.ToEntity(tc.args.d)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
