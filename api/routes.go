package api

import (
	"database/sql"
	"encoding/json"
	"flag"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

var Router = mux.NewRouter()

func listAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := db.ListAllCustomers(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := db.Customer{
		Id: id,
	}
	err = c.GetCustomer(db.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Customer not found"})
		default:
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	utils.ResponseJSON(w, http.StatusOK, c)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var c db.Customer
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload: " + err.Error()})
	}
	defer r.Body.Close()

	err = c.CreateCustomer(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, c)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	var c db.Customer
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&c)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	defer r.Body.Close()

	c.Id = id
	err = c.UpdateCustomer(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["customerId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid customer ID"})
		return
	}

	c := db.Customer{
		Id: id,
	}
	err = c.DeleteCustomer(db.DB)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func addPicture(w http.ResponseWriter, r *http.Request) {

	var p db.PicturePath
	imageName, err := utils.FileUpload(r)
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid data"})
		return
	}

	p.Path = imageName
	err = p.AddPicture(db.DB)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, p)
}

func getPicturePath(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["pictureId"])

	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid picture ID"})
		return
	}

	p := db.PicturePath{
		Id: id,
	}
	err = p.GetPicturePath(db.DB)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Picture not found"})
		default:
			utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	utils.ResponseJSON(w, http.StatusOK, p)
}

func InitRouter() {
	customers := Router.PathPrefix("/customers").Subrouter() // Customer subroute for the API

	customers.HandleFunc("/all", listAllCustomers).Methods("GET")
	customers.HandleFunc("/{customerId:[0-9]+}", getCustomer).Methods("GET")
	customers.HandleFunc("/create", createCustomer).Methods("POST")
	customers.HandleFunc("/{customerId:[0-9]+}", updateCustomer).Methods("PUT")
	customers.HandleFunc("/{customerId:[0-9]+}", deleteCustomer).Methods("DELETE")
	customers.HandleFunc("/picture/{pictureId}", getPicturePath).Methods("GET")
	customers.HandleFunc("/picture", addPicture).Methods("POST")

	var dir string
	flag.StringVar(&dir, "dir", "./img/", "Directory to serve the images")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
}
