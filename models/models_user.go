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
	Username string
	Password []byte
}

func (u *User) CreateUser(db *sql.DB) error {
	passwdHash, err := bcrypt.GenerateFromPassword(u.Password, 14)
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

func (u *User) LoginUser(db *sql.DB) (id int, err error) {
	var passwd []byte
	err = db.QueryRow(`
		SELECT id, passwd FROM users
		WHERE username = $1
		`, u.Username).Scan(id, passwd)
	err = bcrypt.CompareHashAndPassword(passwd, u.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("Invalid credentials")
	} else if err != nil {
		return 0, err
	}
	return id, nil
}
