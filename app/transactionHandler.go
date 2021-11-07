package app

import (
	"encoding/json"
	"encoding/xml"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	"github.com/barnettt/banking/service"
	"github.com/barnettt/banking/util"
	"github.com/gorilla/mux"
	"net/http"
)

type TransactionHandler struct {
	TransactionService service.TransactionService
}

func (transactionHandler *TransactionHandler) SaveTransaction(writer http.ResponseWriter, request *http.Request) {

	contType := request.Header.Get("Content-Type")
	var contentType bool
	if contType == util.ContentTypeXml {
		contentType = true
	}
	inRequest, err := getTransactionRequest(request, contentType)
	if err != nil {
		logger.Error(err.Error())
		appError := exceptions.NewPayloadParseError("Error parsing transaction data")
		transactionHandler.returnResponse(writer, appError, contentType, nil)
	}
	account, appErr := transactionHandler.TransactionService.NewTransaction(inRequest)
	if appErr != nil {
		transactionHandler.returnResponse(writer, appErr, contentType, nil)
		return
	}
	transactionHandler.returnResponse(writer, appErr, contentType, account)
}

func getTransactionRequest(request *http.Request, contentType bool) (*dto.TransactionRequest, error) {
	var transactionRequest *dto.TransactionRequest
	vars := mux.Vars(request)
	var err error
	if contentType {

		err = xml.NewDecoder(request.Body).Decode(&transactionRequest)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	} else {
		err = json.NewDecoder(request.Body).Decode(&transactionRequest)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}
	transactionRequest.AccountId = vars["id"]
	return transactionRequest, err
}

func (transactionHandler *TransactionHandler) returnResponse(writer http.ResponseWriter,
	error *exceptions.AppError, contentType bool, transaction *dto.TransactionResponse) {

	if error != nil {
		if contentType {
			WriteResponse(writer, error.Code, error.AsMessage(), util.ContentTypeXml)
		} else {
			WriteResponse(writer, error.Code, error.AsMessage(), util.ContentTypeJson)
		}
		return
	}
	if contentType {
		// set xml content type on the writer
		WriteResponse(writer, http.StatusOK, transaction, util.ContentTypeXml)
	} else {
		// encode the customers in json format
		WriteResponse(writer, http.StatusOK, transaction, util.ContentTypeJson)
	}
}
