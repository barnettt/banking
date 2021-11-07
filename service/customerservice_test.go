package service

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	domain2 "github.com/barnettt/banking/mock/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*  testing
GetAllCustomers() ([]dto.CustomerResponse, *exceptions.AppError)
GetCustomer(id string) (*dto.CustomerResponse, *exceptions.AppError)
GetCustomersByStatus(status string) ([]dto.CustomerResponse, *exceptions.AppError)
*/

var controller *gomock.Controller
var customerRepository *domain2.MockCustomerRepository
var customerService CustomerService

func customerServiceTestSetup(t *testing.T) func() {
	controller = gomock.NewController(t)
	defer controller.Finish()
	customerRepository = domain2.NewMockCustomerRepository(controller)
	customerService = NewCustomerService(customerRepository)

	return func() {
		customerService = nil
		defer controller.Finish()
	}
}

func getCustomerResponseSlice() []dto.CustomerResponse {
	return []dto.CustomerResponse{
		{"1001", "Ayyub", "Luton", "LT01 8BH", "24/05/1994", "active"},
		{"1002", "Umayamah", "London", "SE6 7TH", "09/02/1997", "active"}}
}

func getCustomerSlice() []domain.Customer {
	return []domain.Customer{
		{"1001", "Ayyub", "Luton", "LT01 8BH", "24/05/1994", "1"},
		{"1002", "Umayamah", "London", "SE6 7TH", "09/02/1997", "1"}}
}

func getSingleCustomerResponse() *dto.CustomerResponse {
	return &dto.CustomerResponse{
		Id: "1001", Name: "Ayyub", City: "Luton", Postcode: "LT01 8BH", DateOfBirth: "24/05/1994", Status: "active"}
}

func getSingleCustomer() *domain.Customer {
	return &domain.Customer{
		Id: "1001", Name: "Ayyub", City: "Luton", Postcode: "LT01 8BH", DateOfBirth: "24/05/1994", Status: "1"}
}

func Test_should_get_all_customers_returning_customer_slice(t *testing.T) {

	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customers := getCustomerSlice()
	responseCustomers := getCustomerResponseSlice()
	customerRepository.EXPECT().FindAll().Return(customers, nil)
	// act
	custs, _ := customerService.GetAllCustomers()
	// assert
	assert.Equal(t, cap(responseCustomers), cap(custs), "Failed to find all customers ")
	assert.Equal(t, responseCustomers, custs, "Failed to find all customers ")
}
func Test_should_get_customers_by_id_returning_a_customer(t *testing.T) {

	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customer := getSingleCustomer()
	expectedCustomer := getSingleCustomerResponse()
	customerRepository.EXPECT().FindById("5678").Return(customer, nil)
	// act
	cust, _ := customerService.GetCustomer("5678")
	// assert
	assert.Equal(t, expectedCustomer, cust, "Failed to find all customers ")
}

/*test shows a pass in console, but for some reason it must be exiting and showing as a non test?? bemused*/
func Test_should_get_customers_by_status_returning_customer(t *testing.T) {
	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customers := getCustomerSlice()
	expectedCustomers := getCustomerResponseSlice()
	customerRepository.EXPECT().FindByStatus("active").Return(customers, nil)
	// act
	cust, _ := customerService.GetCustomersByStatus("active")
	// assert
	assert.Equal(t, expectedCustomers, cust, "Failed to find all customers ")
	assert.Equal(t, "active", cust[0].Status, "Failed customer 0 status not matching")
	assert.Equal(t, "active", cust[1].Status, "Failed customer 1 status not matching")
}

func Test_should_throw_error_when_get_all_customers_is_called(t *testing.T) {
	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customerRepository.EXPECT().FindAll().Return(nil, exceptions.NewDatabaseError("Unexpected database error"))
	// act
	_, err := customerService.GetAllCustomers()
	// assert
	assert.Equal(t, "Unexpected database error", err.Message, "Failed to find all customers ")
}

func Test_should_throw_error_when_get_customer_is_called(t *testing.T) {
	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customerRepository.EXPECT().FindById("5678").Return(nil, exceptions.NewDatabaseError("Unexpected database error"))
	// act
	_, err := customerService.GetCustomer("5678")
	// assert
	assert.Equal(t, "Unexpected database error", err.Message, "Failed to find a customer ")
}

func Test_should_throw_error_when_get_customer_by_status_is_called(t *testing.T) {

	// Arrange
	tearDown := customerServiceTestSetup(t)
	defer tearDown()
	customerRepository.EXPECT().FindByStatus("active").Return(nil, exceptions.NewDatabaseError("Unexpected database error"))
	// act
	_, err := customerService.GetCustomersByStatus("active")
	// assert
	assert.Equal(t, "Unexpected database error", err.Message, "Failed to find by status ")
}
