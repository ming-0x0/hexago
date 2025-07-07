package customer

import (
	"github.com/ming-0x0/hexago/internal/modules/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/modules/customer/port/repository/entity"
)

type CustomerAdapterInterface interface {
	ToDomain(ec *entity.Customer) (*customer.Customer, error)
	ToEntity(ec *customer.Customer) (*entity.Customer, error)
	ToDomains(ecs []*entity.Customer) ([]*customer.Customer, error)
	ToEntities(ecs []*customer.Customer) ([]*entity.Customer, error)
}
