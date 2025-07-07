package repository

import (
	"github.com/ming-0x0/hexago/internal/modules/customer/port/repository/customer"
)

type RepositoryProvider interface {
	CustomerRepository() customer.CustomerRepositoryInterface
}

func NewRepositoryProvider() RepositoryProvider {
	return &repositoryProvider{}
}

type repositoryProvider struct{}

func (rp *repositoryProvider) CustomerRepository() customer.CustomerRepositoryInterface {
	return &customer.CustomerRepository{}
}
