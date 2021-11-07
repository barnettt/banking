package domain

import (
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
)

type Customer struct {
	// set the json output tags in struct
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Postcode    string
	DateOfBirth string `db:"date_of_birth"`
	Status      string
}

func (customer Customer) getStatus() string {
	status := "active"
	if customer.Status == "0" {
		status = "inactive"
	}
	return status
}
func (customer Customer) ToDto() dto.CustomerResponse {

	return dto.CustomerResponse{
		Id:          customer.Id,
		Name:        customer.Name,
		City:        customer.City,
		Postcode:    customer.Postcode,
		DateOfBirth: customer.DateOfBirth,
		Status:      customer.getStatus(),
	}
}

//go:generate mockgen -destination=../mock/domain/mockCustomerRepository.go -package=domain github.com/barnettt/banking/domain CustomerRepository
type CustomerRepository interface {
	FindAll() ([]Customer, *exceptions.AppError)
	FindById(string) (*Customer, *exceptions.AppError)
	FindByStatus(string) ([]Customer, *exceptions.AppError)
}
