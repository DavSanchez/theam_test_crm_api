package main

import (
	"log"
	"net/http"
	"theam.io/jdavidsanchez/test_crm_api/api"
)

func main() {
	api.InitRouter()

	log.Println("Starting server on :4000")

	err := http.ListenAndServe(":4000", api.Router)
	log.Fatal(err)
}
