package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"theam.io/jdavidsanchez/test_crm_api/api"
	"theam.io/jdavidsanchez/test_crm_api/db"
)

func init() {
	db.InitDB()
	api.InitRouter()
}

func main() {
	port := os.Getenv("PORT")
	log.Printf("Starting server on :%s", port)

	server := &http.Server{
		Handler:      api.Router,
		Addr:         ":" + port,
		// Adding timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()
	db.DB.Close()
	log.Fatal(err)
}
