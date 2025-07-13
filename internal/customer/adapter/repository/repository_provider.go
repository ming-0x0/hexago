package repository

import (
	bannerAdapter "github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/banner"
	customerAdapter "github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/customer"
	"github.com/ming-0x0/hexago/internal/customer/port/repository/mysql/banner"
	"github.com/ming-0x0/hexago/internal/customer/port/repository/mysql/customer"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RepositoryProvider struct {
	db              *gorm.DB
	logger          *logrus.Logger
	customerAdapter customer.CustomerRepositoryAdapterInterface
	bannerAdapter   banner.BannerRepositoryAdapterInterface
}

func NewRepositoryProvider(
	db *gorm.DB,
	logger *logrus.Logger,
	customerAdapter customer.CustomerRepositoryAdapterInterface,
	bannerAdapter banner.BannerRepositoryAdapterInterface,
) *RepositoryProvider {
	return &RepositoryProvider{
		db:              db,
		logger:          logger,
		customerAdapter: customerAdapter,
		bannerAdapter:   bannerAdapter,
	}
}

func (rp *RepositoryProvider) CustomerRepository() customer.CustomerRepositoryInterface {
	return customerAdapter.New(rp.db, rp.logger, rp.customerAdapter)
}

func (rp *RepositoryProvider) BannerRepository() banner.BannerRepositoryInterface {
	return bannerAdapter.New(rp.db, rp.logger, rp.bannerAdapter)
}
