package repository

import (
	"github.com/ming-0x0/hexago/internal/customer/port/repository/mysql/customer"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryProviderInterface interface {
	CustomerRepository(
		db *gorm.DB,
		logger *logrus.Logger,
		adapter customer.CustomerRepositoryAdapterInterface,
	) customer.CustomerRepositoryInterface
}
