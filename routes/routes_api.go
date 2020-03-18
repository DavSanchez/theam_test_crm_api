package routes

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/auth"
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
	// users.HandleFunc("/logout", logoutUser).Methods("POST")

	// Static files (customer pictures)
	var dir string
	flag.StringVar(&dir, "dir", "./"+utils.PathToImagesDir+"/", "Directory to serve the images")
	pictures := Router.PathPrefix("/static/").Subrouter()
	pictures.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Registering JWT middleware, Do It Yourself Style!
	customers.Use(auth.ValidateToken)
	pictures.Use(auth.ValidateToken)
}
