package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	db_name = "api"
	db_user = "docker"
	db_pass = "docker"
)

func InitDB() {
	db, err := sql.Open("postgres",fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db_user, db_pass, db_name))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT version()")
	if err != nil {
		log.Fatal(err)
	}

	var result string
	res.Next()
	err = res.Scan(&result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to database: %q", result)
	res.Close()
}
