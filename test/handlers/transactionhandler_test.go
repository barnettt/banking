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

var mockTransactionService *service.MockTransactionService

func transactionTestSetup(t *testing.T) func() {
	// Arrange
	controller = gomock.NewController(t)
	mockTransactionService = service.NewMockTransactionService(controller)
	// account handler
	transactionHandler := app.TransactionHandler{TransactionService: mockTransactionService}

	// create router and routes
	router = mux.NewRouter()
	router.HandleFunc("/customers/{customer_id:[0-9]+}/accounts/{id:[0-9]+}", transactionHandler.SaveTransaction)

	return func() {
		router = nil
		defer controller.Finish()
	}
}

func Test_should_save_new_transaction_with_http_200(t *testing.T) {

	tearDown := transactionTestSetup(t)
	defer tearDown()
	newTransactionRequest := &dto.TransactionRequest{
		AccountId:       "1234",
		Amount:          1000.00,
		TransactionType: "deposit",
	}
	transactionResponse := dto.GetTransactionResponse("225567")
	mockTransactionService.EXPECT().NewTransaction(newTransactionRequest).Return(&transactionResponse, nil)

	// create http request,
	request, _ := http.NewRequest(http.MethodPost, "/customers/5678/accounts/1234", getTransactionAsJson(newTransactionRequest))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code, "Failed to save transaction")
	assert.Contains(t, recorder.Body.String(), "225567", "Failed to create new  transaction")
}

func Test_should_throw_internalServerError_when_save_new_transaction_called(t *testing.T) {

	tearDown := transactionTestSetup(t)
	defer tearDown()
	newTransactionRequest := &dto.TransactionRequest{
		AccountId:       "1234",
		Amount:          1000.00,
		TransactionType: "deposit",
	}
	// transactionResponse := dto.GetTransactionResponse("225567")
	appError := exceptions.NewDatabaseError("Unexpected database error")
	mockTransactionService.EXPECT().NewTransaction(newTransactionRequest).Return(nil, appError)

	// create http request,
	request, _ := http.NewRequest(http.MethodPost, "/customers/5678/accounts/1234", getTransactionAsJson(newTransactionRequest))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Failed to save transaction")
	assert.Contains(t, recorder.Body.String(), "Unexpected database error", "Failed to save new  transaction")
}

func getTransactionAsJson(transaction *dto.TransactionRequest) io.Reader {
	buf := new(bytes.Buffer)
	err := json2.NewEncoder(buf).Encode(transaction)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(buf.String())
}
