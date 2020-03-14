package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var ( /* constants */
	db_name = os.Getenv("DB_NAME")
	db_user = os.Getenv("DB_USER")
	db_pass = os.Getenv("DB_PASS")
	db_host = os.Getenv("DB_HOST")
)

// TODO: Need a Database object for the routes to access it!!

func InitDB() {
	time.Sleep(5 * time.Second) // Wait for database to set up (Docker, just in case!)
	db_params_str := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", db_name, db_user, db_pass, db_host)
	db, err := sql.Open("postgres", db_params_str)

	CheckErr(err)
	defer db.Close()

	err = db.Ping()

	CheckErr(err)
	log.Print("Connected to database")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			username VARCHAR(32) UNIQUE NOT NULL,
			password VARCHAR(64) NOT NULL
		)`)
	CheckErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS customers(
			id SERIAL PRIMARY KEY,
			name VARCHAR(32) NOT NULL,
			surname VARCHAR(32) NOT NULL,
			photo_url TEXT,
			last_modified_by_user_id INT REFERENCES users(id)
		)`)
	CheckErr(err)

	err = InsertUserIfNotExists(db, "Admin", "secret123")
	CheckErr(err)

}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
