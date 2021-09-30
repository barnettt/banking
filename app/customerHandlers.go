package app

import (
	"encoding/json"
	"encoding/xml"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/service"
	"github.com/gorilla/mux"
	"net/http"
)

const contentTypeJson string = "application/json"
const contentTypeXml string = "application/xml"

type CustomerHandler struct {
	service service.CustomerService
}

func (customerHandler *CustomerHandler) GetCustomersByStatus(writer http.ResponseWriter, request *http.Request) {
	print("Called Get Customers by status \n")
	vars := mux.Vars(request)
	status := vars["status"]
	customers, error := customerHandler.service.GetCustomersByStatus(status)
	contentType := request.Header.Get("Content-type") == contentTypeXml
	customerHandler.returnResponse(writer, error, contentType, customers)
}

func (customerHandler *CustomerHandler) getAllCustomers(writer http.ResponseWriter, request *http.Request) {
	print("Called Get all Customers\n")

	status, ok := request.URL.Query()["status"]
	requiredStatus := "1"
	if ok {
		if status[0] == "active" {
			requiredStatus = "1"
		} else {
			requiredStatus = "0"
		}
		customers, error := customerHandler.service.GetCustomersByStatus(requiredStatus)
		contentType := request.Header.Get("Content-type") == contentTypeXml
		customerHandler.returnResponse(writer, error, contentType, customers)
		return
	}
	customers, error := customerHandler.service.GetAllCustomers()
	contentType := request.Header.Get("Content-type") == contentTypeXml
	customerHandler.returnResponse(writer, error, contentType, customers)

}

func (customerHandler *CustomerHandler) getCustomer(writer http.ResponseWriter, request *http.Request) {
	print("Called get a customer  \n")
	vars := mux.Vars(request)
	customer, error := customerHandler.service.GetCustomer(vars["id"])
	contentType := request.Header.Get("Content-Type") == contentTypeXml
	if error != nil {
		if contentType {
			writeResponse(writer, error.Code, error.AsMessage(), contentTypeXml)
		} else {
			writeResponse(writer, error.Code, error.AsMessage(), contentTypeJson)
		}
		return
	}
	if contentType {
		// set xml content type on the writer
		writeResponse(writer, http.StatusOK, customer, contentTypeXml)
	} else {
		// encode the customers in json format
		writeResponse(writer, http.StatusOK, customer, contentTypeJson)
	}

}
func (customerHandler *CustomerHandler) returnResponse(writer http.ResponseWriter, error *exceptions.AppError, contentType bool, customers []dto.CustomerResponse) {
	if error != nil {
		if contentType {
			writeResponse(writer, error.Code, error.AsMessage(), contentTypeXml)
		} else {
			writeResponse(writer, error.Code, error.AsMessage(), contentTypeJson)
		}
		return
	}
	if contentType {
		// set xml content type on the writer
		writeResponse(writer, http.StatusOK, customers, contentTypeXml)
	} else {
		// encode the customers in json format
		writeResponse(writer, http.StatusOK, customers, contentTypeJson)
	}
}
func writeResponse(writer http.ResponseWriter, code int, data interface{}, contentType string) {
	writer.Header().Add("Content-Type", contentType)
	writer.WriteHeader(code)
	if contentType == contentTypeXml {
		error := xml.NewEncoder(writer).Encode(data)
		if error != nil {
			panic(error)
		}
		return
	}
	error := json.NewEncoder(writer).Encode(data)
	if error != nil {
		panic(error)
	}
}
