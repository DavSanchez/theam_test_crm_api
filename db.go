package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func initDB() {
	db, err := sql.Open("postgres","") // TODO

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT version()")
	if err != nil {
		log.Fatal(err)
	}

	var result string
	err = res.Scan(&result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to database: %q", result)
	res.Close()
}
