package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/db"
)

var Router = mux.NewRouter()

func listAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := db.ListAllCustomers(db.DB)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	responseJSON(w, http.StatusOK, customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := db.Customer{
		Id: id,
	}
	err = c.GetCustomer(db.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			responseJSON(w, http.StatusNotFound, map[string]string{"error": "Customer not found"})
		default:
			responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	responseJSON(w, http.StatusOK, c)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var c db.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	defer r.Body.Close()

	err := c.CreateCustomer(db.DB)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusCreated, c)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	var c db.Customer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	defer r.Body.Close()

	c.Id = id
	err = c.UpdateCustomer(db.DB)
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		responseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := db.Customer{
		Id: id,
	}
	err = c.DeleteCustomer(db.DB)

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	responseJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func InitRouter() {
	customers := Router.PathPrefix("/customers").Subrouter() // Customer subroute for the API

	customers.HandleFunc("/all", listAllCustomers).Methods("GET")
	customers.HandleFunc("/{customerId:[0-9]+}", getCustomer).Methods("GET")
	customers.HandleFunc("/create", createCustomer).Methods("POST")
	customers.HandleFunc("/{customerId:[0-9]+}", updateCustomer).Methods("PUT")
	customers.HandleFunc("/{customerId:[0-9]+}", deleteCustomer).Methods("DELETE")
}

func responseJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}
