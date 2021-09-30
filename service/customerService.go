package service

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
)

type CustomerService interface {
	GetAllCustomers() ([]dto.CustomerResponse, *exceptions.AppError)
	GetCustomer(id string) (*dto.CustomerResponse, *exceptions.AppError)
	GetCustomersByStatus(status string) ([]dto.CustomerResponse, *exceptions.AppError)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func (defaultCustomerService DefaultCustomerService) GetAllCustomers() ([]dto.CustomerResponse, *exceptions.AppError) {
	customers, err := defaultCustomerService.repository.FindAll()
	if err != nil {
		return nil, err
	}
	customerResponses := make([]dto.CustomerResponse, 0)
	for i, customer := range customers {
		customerResponses = append(customerResponses, domain.Customer.ToDto(customer))
		i = i + 1 // not used variable but need to take it in the range loop
	}
	return customerResponses, nil
}

func (defaultCustomerService DefaultCustomerService) GetCustomersByStatus(status string) ([]dto.CustomerResponse, *exceptions.AppError) {
	customers, err := defaultCustomerService.repository.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	customerResponses := make([]dto.CustomerResponse, 0)
	for i, customer := range customers {
		customerResponses = append(customerResponses, domain.Customer.ToDto(customer))
		print(i)
	}
	return customerResponses, nil
}

func (defaultCustomerService DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *exceptions.AppError) {
	customer, err := defaultCustomerService.repository.FindById(id)
	if err != nil {
		return nil, err
	}
	customerResponse := customer.ToDto()
	return &customerResponse, nil

}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repo}
}
