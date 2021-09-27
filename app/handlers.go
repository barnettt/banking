package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

const contentTypeJson string = "application/json"
const contentTypeXml string = "application/xml"

type Customer struct {
	// set the json output tags in struct
	Name     string `json:"full_name" xml:"name"`
	City     string `json:"city" xml:"city"`
	Postcode string `json:"post_code"  xml:"postcode"`
}

func greet(writer http.ResponseWriter, request *http.Request) {
	log.Fatal(fmt.Fprintf(writer, "Hello world!!"))
}

func getAllCustomers(writer http.ResponseWriter, request *http.Request) {
	print("Called Get all Customers\n")
	// creating slice of customers with the Customer struct
	customers := []Customer{
		{"Ayyub", "Luton", "LT01 8BH"}, {"Umayamah", "London", "SE6 7TH"},
		{"Sumayah", "London", "BR3 2QQ"},
	}

	contentType := request.Header.Get("Content-type") == contentTypeXml
	// set json content type on the writer - default
	writer.Header().Add("Content-Type", contentTypeJson)
	if contentType {
		// set xml content type on the writer
		writer.Header().Add("Content-Type", contentTypeXml)
		// encode the customers in xml format
		xml.NewEncoder(writer).Encode(customers)
	} else {
		// encode the customers in json format
		json.NewEncoder(writer).Encode(customers)
	}

}
