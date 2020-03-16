package db

import (
	"database/sql"
	"log"
	"os"
	"path"

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
			password BYTEA NOT NULL
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS pictures (
			id SERIAL PRIMARY KEY,
			path TEXT UNIQUE NOT NULL
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			name VARCHAR(32) NOT NULL,
			surname VARCHAR(32) NOT NULL,
			pictureId INT REFERENCES pictures(id),
			lastModifiedByUserId INT REFERENCES users(id)
		)`)
	utils.CheckErr(err)

	initialUser := User{
		Username: "Admin",
		Password: "Secret123",
	}
	noPicturePlaceholder := PicturePath{
		Id: 1,
		Path: path.Join(utils.PathFileServer, "noPicturePlaceholder.jpg"),
	}

	err = initialUser.InsertUserIfNotExists(DB)
	utils.CheckErr(err)
	err = noPicturePlaceholder.AddPlaceholderPicture(DB)
	utils.CheckErr(err)
}
