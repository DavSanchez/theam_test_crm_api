package db

import "database/sql"

// Customer
type Customer struct {
	Id int 						`json:"id"`
	Name string					`json:"name"`
	Surname string				`json:"surname"`
	PhotoURL sql.NullString		`json:"photoUrl"`
	LastModifiedByUserId int	`json:"lastModifiedByUserId"`
	// TODO: What about CreatedByUserId? Should we include time of creation/modification?
}

// User
type User struct {
	Id int
	Username string
	Password string
}