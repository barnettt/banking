package handlers

import (
	"bytes"
	json2 "encoding/json"
	"github.com/barnettt/banking/app"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/mock/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//"/customers/{id:[0-9]+}/accounts"

var mockAccountService *service.MockAccountService

func accountTestSetup(t *testing.T) func() {
	// Arrange
	controller = gomock.NewController(t)
	mockAccountService = service.NewMockAccountService(controller)
	// account handler
	accountHandler := app.AccountHandler{AccountService: mockAccountService}

	// create router and routes
	router = mux.NewRouter()
	router.HandleFunc("/accounts", accountHandler.SaveAccount)
	router.HandleFunc("/customers/{id:[0-9]+}/accounts", accountHandler.SaveAccount)

	return func() {
		router = nil
		defer controller.Finish()
	}
}

func Test_should_create_new_account_for_customer_with_http201(t *testing.T) {

	// Arrange
	tearDown := accountTestSetup(t)
	defer tearDown()
	// set expectations
	newAccount := &dto.AccountRequest{CustomerId: "5678", OpeningDate: "05/11/2021", AccountType: "checking", Amount: 6800.00, Status: "1"}
	newAccountResponse := dto.GetAccountResponse("1234")
	mockAccountService.EXPECT().Save(newAccount).Return(&newAccountResponse, nil)
	// create http request,
	request, _ := http.NewRequest(http.MethodPost, "/customers/5678/accounts", getAsJson(newAccount))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code, "Failed new account save incorrect error code")

}

func Test_should_throw_InternalServerError_when_new_account_for_customer_called(t *testing.T) {

	// Arrange
	tearDown := accountTestSetup(t)
	defer tearDown()
	// set expectations
	newAccount := &dto.AccountRequest{CustomerId: "5678", OpeningDate: "05/11/2021", AccountType: "checking", Amount: 6800.00, Status: "1"}
	//newAccountResponse := dto.GetAccountResponse("1234")
	err := &exceptions.AppError{Code: http.StatusInternalServerError, Message: "Unable to decode payload decoder error "}
	mockAccountService.EXPECT().Save(newAccount).Return(nil, err)
	// create http request,
	request, _ := http.NewRequest(http.MethodPost, "/customers/5678/accounts", getAsJson(newAccount))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Failed to process new account")
	assert.Equal(t, err.Message, "Unable to decode payload decoder error ")

}

func getAsJson(account *dto.AccountRequest) io.Reader {
	buf := new(bytes.Buffer)
	err := json2.NewEncoder(buf).Encode(account)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(buf.String())
}
