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
	db_params_str := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db_user, db_pass, db_name)
	db, err := sql.Open("postgres", db_params_str)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
}

func 