package customer

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/customer/domain/service_type"
	"github.com/ming-0x0/hexago/internal/customer/domain/status"
	"github.com/ming-0x0/hexago/internal/shared/domain/email"
)

type CustomerRepositoryAdapter struct{}

func (a *CustomerRepositoryAdapter) ToDomain(e *entity.Customer) (*customer.Customer, error) {
	email, err := email.New(e.Email)
	if err != nil {
		return nil, err
	}

	serviceType, err := service_type.New(e.ServiceType)
	if err != nil {
		return nil, err
	}

	status, err := status.New(e.Status)
	if err != nil {
		return nil, err
	}

	return customer.FromRepository(
		customer.ID(e.ID),
		e.CustomerName,
		*email,
		e.PhoneNumber,
		e.CompanyName,
		e.Message,
		e.Note,
		*serviceType,
		*status,
	)
}

func (a *CustomerRepositoryAdapter) ToEntity(d *customer.Customer) (*entity.Customer, error) {
	return &entity.Customer{
		ID:           string(d.ID()),
		CustomerName: d.CustomerName(),
		Email:        d.Email().Value(),
		PhoneNumber:  d.PhoneNumber(),
		CompanyName:  d.CompanyName(),
		Message:      d.Message(),
		Note:         d.Note(),
		ServiceType:  d.ServiceType().Value(),
		Status:       d.Status().Value(),
	}, nil
}

func (a *CustomerRepositoryAdapter) ToDomains(es []*entity.Customer) ([]*customer.Customer, error) {
	customers := make([]*customer.Customer, 0, len(es))
	for _, e := range es {
		customer, err := a.ToDomain(e)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (a *CustomerRepositoryAdapter) ToEntities(ds []*customer.Customer) ([]*entity.Customer, error) {
	customers := make([]*entity.Customer, 0, len(ds))
	for _, d := range ds {
		customer, err := a.ToEntity(d)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
