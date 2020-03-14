package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var Router = mux.NewRouter()

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func listAllCustomers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List all customers"))
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customer := params["customer"]

	fmt.Fprintf(w, "Getting customer %s", customer)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a customer"))
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customer := params["customer"]

	fmt.Fprintf(w, "Updating customer %s", customer)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customer := params["customer"]

	fmt.Fprintf(w, "Deleting customer %s", customer)
}

func InitRouter() {
	Router.HandleFunc("/", home) // API root, testing

	customers := Router.PathPrefix("/customers").Subrouter() // Customer subroute for the API

	customers.HandleFunc("/list", listAllCustomers).Methods("GET")
	customers.HandleFunc("/{customer}", getCustomer).Methods("GET")
	customers.HandleFunc("/", createCustomer).Methods("POST")
	customers.HandleFunc("/{customer}", updateCustomer).Methods("PUT")
	customers.HandleFunc("/{customer}", deleteCustomer).Methods("DELETE")
}
