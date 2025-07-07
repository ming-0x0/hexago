package customer

import (
	"github.com/ming-0x0/hexago/internal/modules/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/modules/customer/port/repository/entity"
)

type CustomerAdapter struct{}

func (a *CustomerAdapter) ToDomain(ec *entity.Customer) (*customer.Customer, error) {
	return &customer.Customer{}, nil
}

func (a *CustomerAdapter) ToEntity(ec *customer.Customer) (*entity.Customer, error) {
	return &entity.Customer{}, nil
}

func (a *CustomerAdapter) ToDomains(ecs []*entity.Customer) ([]*customer.Customer, error) {
	return []*customer.Customer{}, nil
}

func (a *CustomerAdapter) ToEntities(ds []*customer.Customer) ([]*entity.Customer, error) {
	return []*entity.Customer{}, nil
}
