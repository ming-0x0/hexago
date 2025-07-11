package customer

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/customer"
)

type CustomerRepositoryAdapterInterface interface {
	ToDomain(e *entity.Customer) (*customer.Customer, error)
	ToEntity(d *customer.Customer) (*entity.Customer, error)
	ToDomains(es []*entity.Customer) ([]*customer.Customer, error)
	ToEntities(ds []*customer.Customer) ([]*entity.Customer, error)
}
