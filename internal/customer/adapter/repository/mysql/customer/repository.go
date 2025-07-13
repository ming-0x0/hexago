package customer

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/customer"
	"github.com/ming-0x0/hexago/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	customerPort "github.com/ming-0x0/hexago/internal/customer/port/repository/mysql/customer"
)

type CustomerRepository struct {
	*repository.Repository[customerPort.CustomerRepositoryAdapterInterface, customer.Customer, entity.Customer]
}

func New(
	db *gorm.DB,
	logger *logrus.Logger,
	adapter customerPort.CustomerRepositoryAdapterInterface,
) *CustomerRepository {
	return &CustomerRepository{
		Repository: repository.NewRepository(db, logger, adapter),
	}
}
