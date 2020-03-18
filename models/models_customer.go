package models

import (
	"database/sql"
	"errors"
)

// Customer
type Customer struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	PictureId            int    `json:"pictureId"`
	CreatedByUserId      int    `json:"createdByUserId"`
	LastModifiedByUserId int    `json:"lastModifiedByUserId"`
}

// Functions for interacting with DB

func (c *Customer) GetCustomer(db *sql.DB) error {
	return db.QueryRow(`
		SELECT customername, surname, pictureId, createdByUserId, lastModifiedByUserId
		FROM customers
		WHERE id = $1
		`, c.Id).Scan(&c.Name, &c.Surname, &c.PictureId, &c.LastModifiedByUserId)
}

func (c *Customer) CreateCustomer(db *sql.DB) error {
	err := db.QueryRow(`
		INSERT INTO customers (customername, surname, pictureId, createdByUserId, lastModifiedByUserId)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id
		`, c.Name, c.Surname, c.PictureId, c.CreatedByUserId).Scan(&c.Id)

	if err != nil {
		return err
	}
	return nil
}

func (c *Customer) UpdateCustomer(db *sql.DB) error {
	res, err := db.Exec(`
		UPDATE customers SET
		customername = $1,
		surname = $2,
		pictureId = $3,
		lastModifiedByUserId = $5
		WHERE id = $6
		`, c.Name, c.Surname, c.PictureId, c.LastModifiedByUserId, c.Id)

	if numRows, _ := res.RowsAffected(); numRows == 0 {
		err = errors.New("No customer was updated")
	}
	return err
}

func (c *Customer) DeleteCustomer(db *sql.DB) error {
	res, err := db.Exec(`
		DELETE FROM customers
		WHERE id = $1
		`, c.Id)

	if numRows, _ := res.RowsAffected(); numRows == 0 {
		err = errors.New("No customer was deleted")
	}

	return err
}

func ListAllCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query(`
		SELECT id, customername, surname, pictureId, createdByUserId, lastModifiedByUserId
		FROM customers`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.Surname, &c.PictureId, &c.CreatedByUserId, &c.LastModifiedByUserId)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}
