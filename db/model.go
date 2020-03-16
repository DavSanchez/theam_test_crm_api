package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"theam.io/jdavidsanchez/test_crm_api/utils"
)

// Customer
type Customer struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	PictureId            int    `json:"pictureId"`
	LastModifiedByUserId int    `json:"lastModifiedByUserId"`
	// TODO: What about CreatedByUserId? Should we include time of creation/modification?
	// pq.NullTime type
}

// User
type User struct {
	Id       int
	Username string
	Password string
}

type PicturePath struct {
	Id   int    `json:"id"`
	Path string `json:"picture"`
}

// // Handling null strings
// type NullInt32 struct {
// 	sql.NullInt32
// }

// func (ni NullInt32) MarshalJSON() ([]byte, error) {
// 	if !ni.Valid {
// 		return []byte("null"), nil
// 	}
// 	return json.Marshal(ni.Int32)
// }

// func (ni *NullInt32) UnmarshalJSON(b []byte) error {
// 	err := json.Unmarshal(b, &ni.Int32)
// 	ni.Valid = (err == nil)
// 	return err
// }

// Functions for interacting with DB

func (c *Customer) GetCustomer(db *sql.DB) error {
	return db.QueryRow(`
		SELECT name, surname, pictureId, lastModifiedByUserId
		FROM customers
		WHERE id = $1
		`, c.Id).Scan(&c.Name, &c.Surname, &c.PictureId, &c.LastModifiedByUserId)
}

func (c *Customer) CreateCustomer(db *sql.DB) error {
	err := db.QueryRow(`
		INSERT INTO customers (name, surname, pictureId, lastModifiedByUserId)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`, c.Name, c.Surname, c.PictureId, c.LastModifiedByUserId).Scan(&c.Id)

	if err != nil {
		return err
	}
	return nil
}

func (c *Customer) UpdateCustomer(db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE customers SET
		name = $1,
		surname = $2,
		pictureId = $3,
		lastModifiedByUserId = $4
		WHERE id = $5
		`, c.Name, c.Surname, c.PictureId, c.LastModifiedByUserId, c.Id)

	return err
}

func (c *Customer) DeleteCustomer(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM customers
		WHERE id = $1
		`, c.Id)

	return err
}

func ListAllCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query(`
		SELECT (id, name, surname, pictureId, lastModifiedByUserId)
		FROM customers`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.Surname, &c.PictureId, &c.LastModifiedByUserId)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

func (u *User) InsertUserIfNotExists(db *sql.DB) error {
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	utils.CheckErr(err)

	_, err = db.Exec(`
		INSERT INTO users (username, password)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`, u.Username, passwdHash)

	if err != nil {
		return err
	}
	fmt.Println("Inserted initial user if it didn't exist ")
	return nil
}

func (p *PicturePath) AddPictureIfNotExists(db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO pictures (path)
		VALUES (1, $1) ON CONFLICT DO NOTHING
		`, p.Path)

	if err != nil {
		return err
	}
	fmt.Println("Inserted placeholder picture if it didn't exist ")
	return nil
}

func (p *PicturePath) AddPicture(db *sql.DB) error {
	err := db.QueryRow(`
		INSERT INTO pictures (path)
		VALUES ($1)
		RETURNING id
		`, p.Path).Scan(&p.Id)

	if err != nil {
		return err
	}
	return nil
}

func (p *PicturePath) GetPicturePath(db *sql.DB) error {
	return db.QueryRow(`
		SELECT path FROM pictures
		WHERE id = $1
		`, p.Id).Scan(&p.Path)
}
