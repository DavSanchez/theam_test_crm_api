package models

import "database/sql"

type PicturePath struct {
	Id   int    `json:"id"`
	Path string `json:"picturePath"`
}

func (p *PicturePath) AddPicture(db *sql.DB) error {
	err := db.QueryRow(`
		INSERT INTO pictures (picturePath)
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
		SELECT picturePath FROM pictures
		WHERE id = $1
		`, p.Id).Scan(&p.Path)
}
