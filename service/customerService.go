package service

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/exceptions"
)

type CustomerService interface {
	GetAllCustomers() ([]domain.Customer, *exceptions.AppError)
	GetCustomer(id string) (*domain.Customer, *exceptions.AppError)
	GetCustomersByStatus(status string) ([]domain.Customer, *exceptions.AppError)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func (defaultCustomerService DefaultCustomerService) GetAllCustomers() ([]domain.Customer, *exceptions.AppError) {
	return defaultCustomerService.repository.FindAll()
}
func (defaultCustomerService DefaultCustomerService) GetCustomersByStatus(status string) ([]domain.Customer, *exceptions.AppError) {
	return defaultCustomerService.repository.FindByStatus(status)
}

func (defaultCustomerService DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *exceptions.AppError) {
	return defaultCustomerService.repository.FindById(id)
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repo}
}
