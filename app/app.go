package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartApp() {

	// create a new multiplexer
	print("creating mux\n ")
	// mux := http.NewServeMux()

	router := mux.NewRouter()
	// define all the routes
	router.HandleFunc("/greet", greet).Methods(http.MethodGet)
	router.HandleFunc("/customers", getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", getCustomer).Methods(http.MethodGet)
	router.HandleFunc("/customers", createCustomer).Methods(http.MethodPost)

	//  start the server using the defaultServMux default multiplexer
	// log any error to fatal
	print("starting listener ..... \n")
	log.Fatal(http.ListenAndServe("localhost:8000", router))

}

func createCustomer(outWriter http.ResponseWriter, inRequest *http.Request) {
	fmt.Fprint(outWriter, "Post Data Received ")
}

func getCustomer(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	fmt.Fprintf(writer, vars["customer_id"])

}
