package app

import (
	"fmt"
	"github.com/barnettt/banking-lib/logger"
	"github.com/barnettt/banking/auth"
	"github.com/barnettt/banking/db"
	"github.com/barnettt/banking/domain"
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
	var dbClient = getDbClient()
	// create a new multiplexer
	// print("creating mux\n ")
	logger.Info("creating mux ")
	// mux := http.NewServeMux()
	// transactionManager := getTransactionManager(dbClient)
	router := mux.NewRouter()
	// Wiring app components
	// handler := CustomerHandler{Service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	customerHandler := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryDb(dbClient))}
	accountHandler := AccountHandler{service.NewAccountService(domain.NewAccountRepositoryDb(dbClient))}
	transactionHandler := TransactionHandler{service.NewTransactionService(domain.NewTransactionRepositoryDb(dbClient), domain.NewAccountRepositoryDb(dbClient), db.NewTxManager(dbClient))}

	// define all the routes

	router.HandleFunc("/customers/{id:[0-9]+}", customerHandler.GetCustomer).Methods(http.MethodGet).Name("GetCustomer")
	router.HandleFunc("/customers?status=active", customerHandler.GetAllCustomers).Methods(http.MethodGet).Name("GetAllActiveCustomer")
	router.HandleFunc("/customers?status=inactive", customerHandler.GetAllCustomers).Methods(http.MethodGet).Name("GetAllInActiveCustomer")
	router.HandleFunc("/customers", customerHandler.GetAllCustomers).Methods(http.MethodGet).Name("GetAllCustomer")
	router.HandleFunc("/customers/{id:[0-9]+}/accounts", accountHandler.SaveAccount).Methods(http.MethodPost).Name("NewAccount")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/accounts/{id:[0-9]+}", transactionHandler.SaveTransaction).Methods(http.MethodPost).Name("NewTransaction")
	authMiddleware := auth.NewAuthorisationMiddleware(domain.AuthorisationRepositoryDB{Client: dbClient})
	router.Use(authMiddleware.AuthorisationHandler())
	//  start the server using the defaultServMux default multiplexer
	// log any error to fatal
	// print("starting listener ..... \n")
	port := os.Getenv("SERVER_PORT")
	host := os.Getenv("SERVER_HOST")
	logger.Info("starting listener ..... on server port : " + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), router))

}

func getTransactionManager(client *sqlx.DB) db.TxManager {
	return db.NewTxManager(client)
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
