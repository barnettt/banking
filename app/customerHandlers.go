package app

import (
	"encoding/json"
	"encoding/xml"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/service"
	"github.com/barnettt/banking/util"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomerHandler struct {
	Service service.CustomerService
}

func (customerHandler *CustomerHandler) GetCustomersByStatus(writer http.ResponseWriter, request *http.Request) {
	print("Called Get Customers by status \n")
	vars := mux.Vars(request)
	status := vars["status"]
	customers, err := customerHandler.Service.GetCustomersByStatus(status)
	contentType := request.Header.Get("Content-type") == util.ContentTypeXml
	customerHandler.returnResponse(writer, err, contentType, customers)
}

func (customerHandler *CustomerHandler) GetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	print("Called Get all Customers\n")

	status, ok := request.URL.Query()["status"]
	requiredStatus := "1"
	if ok {
		if status[0] == "active" {
			requiredStatus = "1"
		} else {
			requiredStatus = "0"
		}
		customers, err := customerHandler.Service.GetCustomersByStatus(requiredStatus)
		contentType := request.Header.Get("Content-type") == util.ContentTypeXml
		customerHandler.returnResponse(writer, err, contentType, customers)
		return
	}
	customers, err := customerHandler.Service.GetAllCustomers()
	contentType := request.Header.Get("Content-type") == util.ContentTypeXml
	customerHandler.returnResponse(writer, err, contentType, customers)

}

func (customerHandler *CustomerHandler) GetCustomer(writer http.ResponseWriter, request *http.Request) {
	print("Called get a customer  \n")
	vars := mux.Vars(request)
	customer, err := customerHandler.Service.GetCustomer(vars["id"])
	contentType := request.Header.Get("Content-Type") == util.ContentTypeXml
	if err != nil {
		if contentType {
			WriteResponse(writer, err.Code, err.AsMessage(), util.ContentTypeXml)
		} else {
			WriteResponse(writer, err.Code, err.AsMessage(), util.ContentTypeJson)
		}
		return
	}
	if contentType {
		// set xml content type on the writer
		WriteResponse(writer, http.StatusOK, customer, util.ContentTypeXml)
	} else {
		// encode the customers in json format
		WriteResponse(writer, http.StatusOK, customer, util.ContentTypeJson)
	}

}
func (customerHandler *CustomerHandler) returnResponse(writer http.ResponseWriter, error *exceptions.AppError, contentType bool, customers []dto.CustomerResponse) {
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
		WriteResponse(writer, http.StatusOK, customers, util.ContentTypeXml)
	} else {
		// encode the customers in json format
		WriteResponse(writer, http.StatusOK, customers, util.ContentTypeJson)
	}
}

func WriteResponse(writer http.ResponseWriter, code int, data interface{}, contentType string) {
	writer.Header().Add("Content-Type", contentType)
	writer.WriteHeader(code)
	if contentType == util.ContentTypeXml {
		err := xml.NewEncoder(writer).Encode(data)
		if err != nil {
			panic(err)
		}
		return
	}
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}
