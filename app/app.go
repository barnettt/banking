package app

import (
	"log"
	"net/http"
)

func StartApp() {
	// define all the routes
	http.HandleFunc("/greet", greet)
	http.HandleFunc("/customers", getAllCustomers)

	//  start the server using the defaultServMux default multiplexer
	// log any error to fatal
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
