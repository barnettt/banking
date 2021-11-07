package service

import (
	trueDomian "github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/mock/domain"
	service2 "github.com/barnettt/banking/service"
	"github.com/golang-sql/civil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var controller *gomock.Controller
var accountRepository *domain.MockAccountRepository
var service service2.AccountService

func accountServiceTestSetup(t *testing.T) func() {
	controller = gomock.NewController(t)
	accountRepository = domain.NewMockAccountRepository(controller)
	service = service2.NewAccountService(accountRepository)

	return func() {
		service = nil
		defer controller.Finish()
	}

}

func Test_should_fail_account_amount_validation_when_creating_new_account(t *testing.T) {

	// Arrange
	// accountRepository := NewMockAccountRepository()
	tearDown := accountServiceTestSetup(t)
	tearDown()
	service := service2.NewAccountService(nil)
	dateTime := civil.DateTimeOf(time.Now())
	request := &dto.AccountRequest{CustomerId: "2456", AccountType: "checking", OpeningDate: dateTime.String(), Amount: 60, Status: "1"}
	// account := getDomainAccount(request)
	// Act
	_, err := service.Save(request)
	// Assert
	assert.Equal(t, "Amount is below the required minimum of 5000 ", err.Message, "Failed validation should fail for amount while saving new account")

}
func Test_should_fail_account_type_validation_when_creating_new_account(t *testing.T) {

	// Arrange
	tearDown := accountServiceTestSetup(t)
	defer tearDown()
	service := service2.NewAccountService(nil)
	dateTime := civil.DateTimeOf(time.Now())
	request := &dto.AccountRequest{CustomerId: "2456", AccountType: "posting", OpeningDate: dateTime.String(), Amount: 6000, Status: "1"}
	// account := getDomainAccount(request)
	// Act
	_, err := service.Save(request)
	// Assert
	assert.Equal(t, "Invalid Account type,  must be one of checking or saving ", err.Message, "Failed validation should fail for account type saving new account")

}

func Test_should_fail_account_creation_when_saving_new_account(t *testing.T) {

	// Arrange
	tearDown := accountServiceTestSetup(t)
	defer tearDown()
	request := &dto.AccountRequest{CustomerId: "2456", AccountType: "checking", OpeningDate: "2021-11-07T15:36:18.559298000", Amount: 6000.00, Status: "1"}
	accountWithOutId := trueDomian.Account{
		CustomerId: "2456", AccountType: "checking",
		OpeningDate: "2021-11-07T15:36:18.559298000",
		Amount:      "6000.000000", Status: "1"}
	accountRepository.EXPECT().Save(accountWithOutId).Return(nil, exceptions.NewDatabaseError("Unknown database error"))
	// Act
	_, err := service.Save(request)
	// Assert
	assert.Equal(t, "Unknown database error", err.Message, "Failed while saving new account to repository")

}

func Test_should_create_new_account_with_http201(t *testing.T) {

	// Arrange
	tearDown := accountServiceTestSetup(t)
	defer tearDown()
	request := &dto.AccountRequest{CustomerId: "2456", AccountType: "checking", OpeningDate: "2021-11-07T15:36:18.559298000", Amount: 6000.00, Status: "1"}
	accountWithId := &trueDomian.Account{
		CustomerId: "2456", AccountType: "checking",
		OpeningDate: "2021-11-07T15:36:18.559298000",
		Amount:      "6000.000000", Status: "1"}

	accountRepository.EXPECT().Save(*accountWithId).Return(&dto.NewAccountResponse{AccountId: "5678"}, nil)
	// Act
	accountResponse, _ := service.Save(request)
	// Assert
	assert.Equal(t, "5678", accountResponse.AccountId, "Failed while saving new account to repository")

}
