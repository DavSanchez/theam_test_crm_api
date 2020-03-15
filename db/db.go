package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	utils.CheckErr(err)

	err = DB.Ping()

	utils.CheckErr(err)
	log.Print("Connected to database")

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(32) UNIQUE NOT NULL,
			password VARCHAR(64) NOT NULL
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS pictures (
			id SERIAL PRIMARY KEY,
			path TEXT
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(32) NOT NULL,
			surname VARCHAR(32) NOT NULL,
			pictureId INT,
			lastModifiedByUserId INT REFERENCES users(id) NOT NULL
		)`)
	utils.CheckErr(err)

	initialUser := User{
		Username: "Admin",
		Password: "Secret123",
	}

	err = initialUser.InsertUserIfNotExists(DB)
	utils.CheckErr(err)

}
