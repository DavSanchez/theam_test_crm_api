package main

import (
	"log"
	"net/http"

	"theam.io/jdavidsanchez/test_crm_api/api"
	"theam.io/jdavidsanchez/test_crm_api/db"
)

func init() {
	db.InitDB()
	api.InitRouter()
}

func main() {
	log.Println("Starting server on :4000")

	err := http.ListenAndServe(":4000", api.Router)
	db.DB.Close()
	log.Fatal(err)
}
