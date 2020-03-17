package api

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/db"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

var Router = mux.NewRouter()

func InitRouter() {
	// Customer subroute for the API
	customers := Router.PathPrefix("/customers").Subrouter()

	customers.HandleFunc("/all", listAllCustomers).Methods("GET")
	customers.HandleFunc("/{customerId:[0-9]+}", getCustomer).Methods("GET")
	customers.HandleFunc("/create", createCustomer).Methods("POST")
	customers.HandleFunc("/{customerId:[0-9]+}", updateCustomer).Methods("PUT")
	customers.HandleFunc("/{customerId:[0-9]+}", deleteCustomer).Methods("DELETE")
	customers.HandleFunc("/picture/{pictureId}", getPicturePath).Methods("GET")
	customers.HandleFunc("/picture", addPicture).Methods("POST")

	// User authentication
	users := Router.PathPrefix("/users").Subrouter()

	users.HandleFunc("/register", registerUser).Methods("POST")
	users.HandleFunc("/login", loginUser).Methods("POST")
	users.HandleFunc("/logout", logoutUser).Methods("POST")

	// Static files (customer pictures)
	var dir string
	flag.StringVar(&dir, "dir", "./"+utils.PathToImagesDir+"/", "Directory to serve the images")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
}

/**************
Customer routes
***************/

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
	id, err := strconv.Atoi(params["customerId"]) // This parameter is always an int (Regex in mux route)

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
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
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
		return
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

/*************
Picture routes
**************/

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

/***************
User auth routes
****************/

func registerUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Register user")
}
func loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Login user")
}
func logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout user")
}
