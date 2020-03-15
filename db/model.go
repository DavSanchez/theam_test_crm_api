package db

import (
	"database/sql"
	"encoding/json"
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

func (c *Customer) GetCustomer(db *sql.DB) error {
	return db.QueryRow(`
		SELECT (name, surname, photoUrl, lastModifiedByUserId)
		FROM customers
		WHERE id = $1
		`, c.Id).Scan(&c.Name, &c.Surname, &c.PhotoURL, &c.LastModifiedByUserId)
}

func (c *Customer) CreateCustomer(db *sql.DB) error {
	err := db.QueryRow(`
		INSERT INTO customers (name, surname, photoUrl, lastModifiedByUserId)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`, c.Name, c.Surname, c.PhotoURL, c.LastModifiedByUserId).Scan(&c.Id)

	if err != nil {
		return err
	}
	return nil
}

func (c *Customer) UpdateCustomer(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE customers
		SET (name=$1, surname=$2, photoUrl=$3, lastModifiedByUserId=$4)
		WHERE id=$4
		`, c.Name, c.Surname, c.PhotoURL, c.LastModifiedByUserId, c.Id)

	return err
}

func (c *Customer) DeleteCustomer(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM customers
		WHERE id=$1
		`, c.Id)

	return err
}

func ListAllCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query(`
		SELECT (id, name, surname, photoUrl, lastModifiedByUserId)
		FROM customers`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.Surname, &c.PhotoURL, &c.LastModifiedByUserId)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

func (u *User) InsertUserIfNotExists(db *sql.DB) error {
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	CheckErr(err)

	_, err = db.Exec(`
		INSERT INTO users(username, password)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`, u.Username, string(passwdHash))

	CheckErr(err)

	fmt.Println("Inserted initial user if it didn't exist ")

	return err
}
