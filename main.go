package main

import (
	"log"
	"net/http"
	"os"

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

	err := http.ListenAndServe(":" + port, api.Router)
	db.DB.Close()
	log.Fatal(err)
}
