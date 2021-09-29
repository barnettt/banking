package app

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartApp() {

	// create a new multiplexer
	print("creating mux\n ")
	// mux := http.NewServeMux()

	router := mux.NewRouter()
	// Wiring app components
	// handler := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	handler := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDb())}

	// define all the routes

	router.HandleFunc("/customers/{id:[0-9]+}", handler.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers?status=active", handler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers?status=inactive", handler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", handler.getAllCustomers).Methods(http.MethodGet)

	//  start the server using the defaultServMux default multiplexer
	// log any error to fatal
	print("starting listener ..... \n")
	log.Fatal(http.ListenAndServe("localhost:8000", router))

}
