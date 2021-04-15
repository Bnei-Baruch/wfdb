package models

import (
	"database/sql"
)

type Kmedia struct {
	ID        int    `json:"id"`
	KmediaID  string `json:"kmedia_id"`
	Date      string `json:"date"`
	FileName  string `json:"file_name"`
	Language  string `json:"language"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Sha1      string `json:"sha1"`
	Pattern   string `json:"pattern"`
	SendID    string `json:"send_id"`
	Source    string `json:"source"`
}

func FindKmFiles(db *sql.DB, key string, value string) ([]Kmedia, error) {
	sqlStatement := `SELECT id, kmedia_id, date, file_name, language, extension, size, sha1, pattern, send_id, source FROM kmedia WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Kmedia{}

	for rows.Next() {
		var a Kmedia
		if err := rows.Scan(&a.ID, &a.KmediaID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Pattern, &a.SendID, &a.Source); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func GetKmFiles(db *sql.DB, start, count int) ([]Kmedia, error) {
	rows, err := db.Query(
		"SELECT id, kmedia_id, date, file_name, language, extension, size, sha1, pattern, send_id, source FROM kmedia ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Kmedia{}

	for rows.Next() {
		var a Kmedia
		if err := rows.Scan(&a.ID, &a.KmediaID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Pattern, &a.SendID, &a.Source); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func (a *Kmedia) GetKmFile(db *sql.DB) error {

	return db.QueryRow("SELECT id, kmedia_id, date, file_name, language, extension, size, sha1, pattern, send_id, source FROM kmedia WHERE kmedia_id = $1",
		a.KmediaID).Scan(&a.ID, &a.KmediaID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Pattern, &a.SendID, &a.Source)
}

func (a *Kmedia) PostKmFile(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO kmedia(kmedia_id, date, file_name, language, extension, size, sha1, pattern, send_id, source) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (kmedia_id) DO UPDATE SET (kmedia_id, date, file_name, language, extension, size, sha1, pattern, send_id, source) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE kmedia.kmedia_id = $1 RETURNING id",
		a.KmediaID, a.Date, a.FileName, a.Language, a.Extension, a.Size, a.Sha1, a.Pattern, &a.SendID, &a.Source).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *Kmedia) DeleteKmFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM kmedia WHERE kmedia_id=$1", a.KmediaID)

	return err
}
