package handlers

import (
	"github.com/barnettt/banking/app"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/mock/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var controller *gomock.Controller
var mockCustomerService *service.MockCustomerService
var customerHandler app.CustomerHandler
var router *mux.Router

func testSetup(t *testing.T) func() {
	// Arrange
	controller = gomock.NewController(t)
	mockCustomerService = service.NewMockCustomerService(controller)
	// wire handler with mock
	customerHandler = app.CustomerHandler{Service: mockCustomerService}
	// define the route to be tested
	router = mux.NewRouter()
	router.HandleFunc("/customers", customerHandler.GetAllCustomers)
	router.HandleFunc("/customers/{id}", customerHandler.GetCustomer)
	router.HandleFunc("/customers?status=active", customerHandler.GetCustomersByStatus)
	router.HandleFunc("/customers?status=inactive", customerHandler.GetCustomersByStatus)
	// defer must always be called at end of method
	return func() {
		router = nil
		defer controller.Finish()
	}
}

func Test_should_return_customer_slice_with_status_200(t *testing.T) {

	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// define the data for mock
	customers := customersWithActiveStatus(true)
	mockCustomerService.EXPECT().GetAllCustomers().Return(customers, nil)
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, recorder.Code, http.StatusOK, "Failed getAllCustomerTest 200, incorrect status")
	assert.NotEmpty(t, recorder.Body, "Failed getAllCustomerTest response body empty")
}

func Test_should_return_error_message_with_status_500(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// set up expect to return error
	mockCustomerService.EXPECT().GetAllCustomers().Return(nil, exceptions.NewDatabaseError("Dummy unexpected error"))
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, recorder.Code, http.StatusInternalServerError, "Failed getAllCustomerTest 500,  incorrect status")
	assert.Contains(t, recorder.Body.String(), "Dummy unexpected error", "Failed getAllCustomerTest response body not contain message")
}

func Test_should_get_http200_when_get_customer_by_id_clled(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// set up expect to return customer
	customer := dto.CustomerResponse{Id: "1001", Name: "Ayyub", City: "Luton", Postcode: "LT01 8BH", DateOfBirth: "24/05/1994", Status: "1"}
	mockCustomerService.EXPECT().GetCustomer("1234").Return(&customer, nil)
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers/1234", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, http.StatusOK, "Failed getCustomerByIdTest incorrect status returned")
	assert.Contains(t, recorder.Body.String(), "Ayyub", "getCustomer, name not returned ")
	// id not returned in body
	// assert.Contains(t,recorder.Body.String(),"1234", "getCustomer, id not returned ")
	assert.Contains(t, recorder.Body.String(), "Luton", "getCustomer, city not returned ")
	assert.Contains(t, recorder.Body.String(), "LT01 8BH", "getCustomer, postcode not returned ")
	assert.Contains(t, recorder.Body.String(), "24/05/1994", "getCustomer,dob not returned ")
	assert.Contains(t, recorder.Body.String(), "1", "getCustomer,status not returned ")
}

func Test_should_throw_404_error_when_get_customer_by_id_called(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// set up expect to return error
	mockCustomerService.EXPECT().GetCustomer("1234").Return(nil, exceptions.NewNotFoundError("dummy not found message"))
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers/1234", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, http.StatusNotFound, "Failed getCustomerByIdTest expected 404, incorrect status returned")
	assert.Contains(t, recorder.Body.String(), "dummy not found message", "Failed getCustomerByIdTest 404, incorrect error message")
}

func Test_should_get_http200_when_get_customers_by_status_called_with_active_status(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// define the data for mock
	customers := customersWithActiveStatus(true)
	mockCustomerService.EXPECT().GetCustomersByStatus("1").Return(customers, nil)
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers?status=active", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, recorder.Code, http.StatusOK, "Failed getCustomerByActiveStatus 200, incorrect status")
	assert.NotEmpty(t, recorder.Body, "Failed getCustomerByActiveStatus 200,  response body empty")

}
func Test_should_get_http200_when_get_customers_by_status_called_with_inactive_status(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// define the data for mock
	customers := customersWithActiveStatus(false)
	mockCustomerService.EXPECT().GetCustomersByStatus("0").Return(customers, nil)
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers?status=inactive", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, recorder.Code, http.StatusOK, "Failed getCustomerByInactiveStatus 200, incorrect status")
	assert.NotEmpty(t, recorder.Body, "Failed getCustomerByInactiveStatus 200,  response body empty")

}

func Test_should_throw_http500_when_get_customers_by_status_called_with_any_status(t *testing.T) {
	// Arrange
	tearDown := testSetup(t)
	defer tearDown()
	// define the data for mock
	mockCustomerService.EXPECT().GetCustomersByStatus("1").Return(nil, exceptions.NewDatabaseError("dummy unexpected error"))
	// create a request
	request, _ := http.NewRequest(http.MethodGet, "/customers?status=active", nil)

	// Act
	// create test recorder
	recorder := httptest.NewRecorder()
	// invoke the route
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, recorder.Code, http.StatusInternalServerError, "Failed getCustomerByInactiveStatus 200, incorrect status")
	assert.Contains(t, recorder.Body.String(), "dummy unexpected error", "Failed getCustomerByInactiveStatus 200,  response body empty")

}

func customersWithActiveStatus(status bool) []dto.CustomerResponse {
	var st string
	if status {
		st = "1"
	} else {
		st = "0"
	}
	customers := []dto.CustomerResponse{
		{"1001", "Ayyub", "Luton", "LT01 8BH", "24/05/1994", st},
		{"1002", "Umayamah", "London", "SE6 7TH", "09/02/1997", st},
		{"1003", "Sumayah", "London", "BR3 2QQ", "10/02/1998", st},
	}
	return customers
}
