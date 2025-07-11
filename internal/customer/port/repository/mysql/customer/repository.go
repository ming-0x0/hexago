package customer

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/shared/repository"
)

type CustomerRepositoryInterface interface {
	repository.RepositoryInterface[CustomerRepositoryAdapterInterface, customer.Customer, entity.Customer]
}
