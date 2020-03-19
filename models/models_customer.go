package models

import (
	"database/sql"
	"errors"
)

// Customer
type Customer struct {
	CustomerOut
	PictureId            int `json:"pictureId"`
	CreatedByUserId      int
	LastModifiedByUserId int
}

type CustomerOut struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	PicturePath        string `json:"picturePath"`
	CreatedByUser      string `json:"createdByUser"`
	LastModifiedByUser string `json:"lastModifiedByUser"`
}

// Functions for interacting with DB

func (c *CustomerOut) GetCustomer(db *sql.DB) error {
	return db.QueryRow(`
		SELECT 
		customername,
		surname,
		(SELECT picturePath FROM pictures WHERE id = pictureId),
		(SELECT username FROM users WHERE id = createdByUserId),
		(SELECT username FROM users WHERE id = lastModifiedByUserId)
		FROM customers
		WHERE id = $1
		`, c.Id).Scan(&c.Name, &c.Surname, &c.PicturePath, &c.CreatedByUser, &c.LastModifiedByUser)
}

func (c *Customer) CreateCustomer(db *sql.DB) error {
	pictureId := 1
	if c.PictureId != 0 {
		pictureId = c.PictureId
	}
	err := db.QueryRow(`
		INSERT INTO customers (
			customername,
			surname, 
			pictureId,
			createdByUserId,
			lastModifiedByUserId
		)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, (SELECT picturePath FROM pictures WHERE id = pictureId),
		(SELECT username FROM users WHERE id = createdByUserId),
		(SELECT username FROM users WHERE id = lastModifiedByUserId)
		`, c.Name, c.Surname, pictureId, c.CreatedByUserId).Scan(
		&c.Id, &c.PicturePath, &c.CreatedByUser, &c.LastModifiedByUser)

	if err != nil {
		return err
	}
	return nil
}

func (c *Customer) UpdateCustomer(db *sql.DB) error {
	pictureId := 1
	if c.PictureId != 0 {
		pictureId = c.PictureId
	}
	err := db.QueryRow(`
		UPDATE customers SET
		customername = COALESCE($1, customername),
		surname = COALESCE($2, surname),
		pictureId = COALESCE($3, pictureId),
		lastModifiedByUserId = COALESCE($4, lastModifiedByUserId)
		WHERE id = $5
		RETURNING id, (SELECT picturePath FROM pictures WHERE id = pictureId),
		(SELECT username FROM users WHERE id = createdByUserId),
		(SELECT username FROM users WHERE id = lastModifiedByUserId)
		`, c.Name, c.Surname, pictureId, c.LastModifiedByUserId, c.Id).Scan(
		&c.Id, &c.PicturePath, &c.CreatedByUser, &c.LastModifiedByUser)

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("No customer was updated")
		}
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

func ListAllCustomers(db *sql.DB) ([]CustomerOut, error) {
	rows, err := db.Query(`
		SELECT id, customername, 
		surname, 
		(SELECT picturePath FROM pictures WHERE id = pictureId),
		(SELECT username FROM users WHERE id = createdByUserId),
		(SELECT username FROM users WHERE id = lastModifiedByUserId)
		FROM customers`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []CustomerOut{}

	for rows.Next() {
		var c CustomerOut
		err := rows.Scan(&c.Id, &c.Name, &c.Surname, &c.PicturePath, &c.CreatedByUser, &c.LastModifiedByUser)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}
