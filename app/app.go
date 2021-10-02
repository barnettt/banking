package app

import (
	"fmt"
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/logger"
	"github.com/barnettt/banking/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
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
	// create db connection pool
	dbClient := getDbClient()
	// create a new multiplexer
	// print("creating mux\n ")
	logger.Info("creating mux ")
	// mux := http.NewServeMux()

	router := mux.NewRouter()
	// Wiring app components
	// handler := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	customerHandler := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDb(dbClient))}
	accountHandler := AccountHandler{service.NewAccountService(domain.NewAccountRepositoryDb(dbClient))}

	// define all the routes

	router.HandleFunc("/customers/{id:[0-9]+}", customerHandler.getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers?status=active", customerHandler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers?status=inactive", customerHandler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", customerHandler.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{id:[0-9]+}/accounts", accountHandler.saveAccount).Methods(http.MethodPost)

	//  start the server using the defaultServMux default multiplexer
	// log any error to fatal
	// print("starting listener ..... \n")
	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	logger.Info("starting listener ..... on server port : " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router))

}

func getDbClient() *sqlx.DB {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbDrivername := os.Getenv("DB_DRIVER_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	client, err := sqlx.Open(fmt.Sprintf("%s", dbDrivername), fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", user, password, dbProtocol, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
