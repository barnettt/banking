package app

import (
	"fmt"
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/logger"
	"github.com/barnettt/banking/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func StartApp() {
	if os.Getenv("SERVER_PORT") == "" ||
		os.Getenv("SERVER_HOST") == "" ||
		os.Getenv("DB_HOST") == "" ||
		os.Getenv("DB_PORT") == "" ||
		os.Getenv("DB_USER") == "" ||
		os.Getenv("DB_PASSWD") == "" ||
		os.Getenv("DB_PROTOCOL") == "" ||
		os.Getenv("DB_NAME") == "" ||
		os.Getenv("DB_DRIVER_NAME") == "" {
		logger.Error("Environment variables are undefined ... ")
		log.Fatal("Environment variables are undefined ... ")
	}
	// create a new multiplexer
	// print("creating mux\n ")
	logger.Info("creating mux ")
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
	// print("starting listener ..... \n")
	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	logger.Info("starting listener ..... on server port : " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router))

}
