package routes

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"theam.io/jdavidsanchez/test_crm_api/auth"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

var Router = mux.NewRouter()
var notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
})

func InitRouter() {
	// As it is an API, handle invalid routes with a JSON-formatted 404 Not Found
	Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
	})

	// Customer subroute for the API
	customers := Router.PathPrefix("/customers").Subrouter()

	customers.HandleFunc("/all", listAllCustomers).Methods("GET")
	customers.HandleFunc("/{customerId:[0-9]+}", getCustomer).Methods("GET")
	customers.HandleFunc("/", createCustomer).Methods("POST")
	customers.HandleFunc("/{customerId:[0-9]+}", updateCustomer).Methods("PUT")
	customers.HandleFunc("/{customerId:[0-9]+}", deleteCustomer).Methods("DELETE")
	customers.HandleFunc("/picture/{pictureId:[0-9]+}", getPicturePath).Methods("GET")
	customers.HandleFunc("/picture", addPicture).Methods("POST")
	// User authentication
	users := Router.PathPrefix("/users").Subrouter()

	users.HandleFunc("/register", registerUser).Methods("POST")
	users.HandleFunc("/login", loginUser).Methods("POST")

	// Static files (customer pictures)
	var dir string
	flag.StringVar(&dir, "images", "./"+utils.PathToImagesDir+"/", "Directory to serve the images")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Register JWT middleware
	customers.Use(auth.ValidateToken)

	var publicDir string
	flag.StringVar(&publicDir, "public", "./public/", "Directory to serve the homepage")
	Router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(publicDir))))

	// As it is an API, handle invalid routes with a JSON-formatted 404 Not Found
	Router.NotFoundHandler = notFoundHandler
	customers.NotFoundHandler = notFoundHandler
	users.NotFoundHandler = notFoundHandler
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}