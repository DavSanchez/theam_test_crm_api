package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Customer
type Customer struct {
	Id                   int            `json:"id"`
	Name                 string         `json:"name"`
	Surname              string         `json:"surname"`
	PhotoURL             sql.NullString `json:"photoUrl"`
	LastModifiedByUserId int            `json:"lastModifiedByUserId"`
	// TODO: What about CreatedByUserId? Should we include time of creation/modification?
	// pq.NullTime type
}

// User
type User struct {
	Id       int
	Username string
	Password string
}

// Handling null strings
type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// Functions for interacting with DB

func ListAllCustomers(db *sql.DB) error {
	return errors.New("Error")
}

func GetCustomer(db *sql.DB) error {
	return errors.New("Error")
}

func CreateCustomer(db *sql.DB) error {
	return errors.New("Error")
}

func UpdateCustomer(db *sql.DB) error {
	return errors.New("Error")
}

func DeleteCustomer(db *sql.DB) error {
	return errors.New("Error")
}

func InsertUserIfNotExists(db *sql.DB, username, password string) error {
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	CheckErr(err)

	_, err = db.Exec(`
		INSERT INTO users(username, password)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`, username, string(passwdHash))
	CheckErr(err)

	fmt.Println("Inserted a new user")

	return err
}
