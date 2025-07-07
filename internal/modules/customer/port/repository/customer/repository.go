package customer

import (
	"github.com/ming-0x0/hexago/internal/modules/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/modules/customer/port/repository/entity"
	"github.com/ming-0x0/hexago/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerRepositoryInterface interface {
	repository.RepositoryInterface[CustomerAdapterInterface, customer.Customer, entity.Customer]
}

type CustomerRepository struct {
	*repository.Repository[CustomerAdapterInterface, customer.Customer, entity.Customer]
}

func New(db *gorm.DB, logger *logrus.Logger, adapter CustomerAdapterInterface) CustomerRepositoryInterface {
	return &CustomerRepository{
		repository.NewRepository(db, logger, adapter).(*repository.Repository[CustomerAdapterInterface, customer.Customer, entity.Customer]),
	}
}
