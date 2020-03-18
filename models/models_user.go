package models

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

// User
type User struct {
	Id       int
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) CreateUser(db *sql.DB) error {
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	utils.CheckErr(err)

	_, err = db.Exec(`
		INSERT INTO users (username, passwd)
		VALUES ($1, $2)
		`, u.Username, passwdHash)

	// err.(*pq.Error) is a type assertion
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" && err.Column == "username" {
			// Unique violation of username field
			return errors.New("Username already in use")
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (u *User) LoginUser(db *sql.DB) error {
	passwd := []byte(u.Password)
	err := db.QueryRow(`
		SELECT id, username, passwd FROM users
		WHERE username = $1
		`, u.Username).Scan(&u.Id, &u.Username, &u.Password)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), passwd)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("Invalid credentials")
	} else if err != nil {
		return err
	}
	return nil
}
