package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var ( /* constants and database global object */
	db_name = os.Getenv("DB_NAME")
	db_user = os.Getenv("DB_USER")
	db_pass = os.Getenv("DB_PASS")
	db_host = os.Getenv("DB_HOST")

	DB *sql.DB
)

func InitDB() {
	var err error

	db_params_str := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		db_name,
		db_user,
		db_pass,
		db_host)

	DB, err = sql.Open("postgres", db_params_str)

	CheckErr(err)

	err = DB.Ping()

	CheckErr(err)
	log.Print("Connected to database")

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(32) UNIQUE NOT NULL,
			password VARCHAR(64) NOT NULL
		)`)
	CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(32) NOT NULL,
			surname VARCHAR(32) NOT NULL,
			photoUrl TEXT,
			lastModifiedByUserId INT REFERENCES users(id)
		)`)
	CheckErr(err)

	initialUser := User{
		Username: "Admin",
		Password: "Secret123",
	}

	err = initialUser.InsertUserIfNotExists(DB)
	CheckErr(err)

}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
