package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

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
