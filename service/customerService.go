package service

import "github.com/barnettt/banking/domain"

type CustomerService interface {
	GetAllCustomers() ([]domain.Customer, error)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func (defaultCustomerService DefaultCustomerService) GetAllCustomers() ([]domain.Customer, error) {
	return defaultCustomerService.repository.FindAll()
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repo}
}
