// model_kmedia.go

package main

import (
	"database/sql"
)

type kmedia struct {
	ID    int			`json:"id"`
	KmediaID  string  	`json:"kmedia_id"`
	Date string			`json:"date"`
	FileName string		`json:"file_name"`
	Language string		`json:"language"`
	Extension string	`json:"extension"`
	Size int64			`json:"size"`
	Sha1 string			`json:"sha1"`
	Pattern  string  	`json:"pattern"`
}

func findKmFiles(db *sql.DB, key string, value string) ([]kmedia, error) {
	sqlStatement := `SELECT * FROM kmedia WHERE `+key+` LIKE '%`+value+`' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []kmedia{}

	for rows.Next() {
		var a kmedia
		if err := rows.Scan(&a.ID, &a.KmediaID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Pattern); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func getKmFiles(db *sql.DB, start, count int) ([]kmedia, error) {
	rows, err := db.Query(
		"SELECT * FROM kmedia ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []kmedia{}

	for rows.Next() {
		var a kmedia
		if err := rows.Scan(&a.ID, &a.KmediaID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Pattern); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func (a *kmedia) getKmFile(db *sql.DB) error {

	return db.QueryRow("SELECT * FROM kmedia WHERE kmedia_id = $1",
		a.KmediaID).Scan(&a)
}

func (a *kmedia) postKmFile(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO kmedia(kmedia_id, date, file_name, language, extension, size, sha1, pattern) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (kmedia_id) DO UPDATE SET (kmedia_id, date, file_name, language, extension, size, sha1, pattern) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE kmedia.kmedia_id = $1 RETURNING id",
		a.KmediaID, a.Date, a.FileName, a.Language, a.Extension, a.Size, a.Sha1, a.Pattern).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *kmedia) deleteKmFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM kmedia WHERE kmedia_id=$1", a.KmediaID)

	return err
}