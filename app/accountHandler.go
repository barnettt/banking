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

type AccountHandler struct {
	accountService service.AccountService
}

func (accountHandler *AccountHandler) saveAccount(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	contentType := request.Header.Get("Content-Type")
	account, derr := getDecodedAccount(contentType, request)
	if derr != nil {
		appError := exceptions.NewPayloadParseError("Unable to decode payload decoder error ")
		WriteResponse(writer, http.StatusInternalServerError, appError.AsMessage(), contentType)
		return
	}
	account.CustomerId = vars["id"]
	acc, err := accountHandler.accountService.Save(account)
	if err != nil {
		appError := exceptions.NewPayloadParseError("Unable save to account")
		WriteResponse(writer, http.StatusInternalServerError, appError.AsMessage(), contentType)
		return
	}
	err = err
	WriteResponse(writer, http.StatusCreated, acc, contentType)
}

func getDecodedAccount(contentType string, request *http.Request) (*dto.AccountRequest, error) {

	var account dto.AccountRequest

	if contentType == util.ContentTypeXml {

		err := xml.NewDecoder(request.Body).Decode(&account)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	} else {
		err := json.NewDecoder(request.Body).Decode(&account)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}
	return &account, nil
}
