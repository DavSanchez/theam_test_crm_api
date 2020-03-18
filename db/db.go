package db

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "github.com/lib/pq"
	"theam.io/jdavidsanchez/test_crm_api/models"
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

	// First time initialization (// FIXME ?)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(64) UNIQUE NOT NULL,
			passwd BYTEA NOT NULL
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS pictures (
			id SERIAL PRIMARY KEY,
			picturePath TEXT UNIQUE NOT NULL
		)`)
	utils.CheckErr(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS customers (
			id SERIAL PRIMARY KEY,
			customername VARCHAR(32) NOT NULL,
			surname VARCHAR(32) NOT NULL,
			lastModifiedByUserId INT NOT NULL,
			pictureId INT REFERENCES pictures(id),
			createdByUserId INT REFERENCES users(id)
		)`)
	utils.CheckErr(err)

	initialUser := models.User{
		Username: "Admin",
		Password: "hunter2",
	}
	noPicturePlaceholder := models.PicturePath{
		Id:   1,
		Path: path.Join(utils.PathFileServer, "noPicturePlaceholder.jpg"),
	}

	err = initialUser.CreateUser(DB)
	utils.CheckErr(err)
	err = noPicturePlaceholder.AddPicture(DB)
	utils.CheckErr(err)
}
