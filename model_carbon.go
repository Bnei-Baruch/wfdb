// model_carbon.go

package main

import (
	"database/sql"
)

type carbon struct {
	ID    int			`json:"id"`
	CarbonID  string  	`json:"carbon_id"`
	SendID  string  	`json:"send_id"`
	Date string			`json:"date"`
	FileName string		`json:"file_name"`
	Language string		`json:"language"`
	Extension string	`json:"extension"`
	Duration float32	`json:"duration"`
	Size int64			`json:"size"`
	Sha1 string			`json:"sha1"`
}

func findCarbonFiles(db *sql.DB, key string, value string) ([]carbon, error) {
	sqlStatement := `SELECT id, carbon_id, send_id, date, file_name, language, extension, size, sha1, duration FROM carbon WHERE `+key+` LIKE '%`+value+`' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []carbon{}

	for rows.Next() {
		var a carbon
		if err := rows.Scan(&a.ID, &a.CarbonID, &a.SendID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Duration); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func getCarbonFiles(db *sql.DB, start, count int) ([]carbon, error) {
	rows, err := db.Query(
		"SELECT id, carbon_id, send_id, date, file_name, language, extension, size, sha1, duration FROM carbon ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []carbon{}

	for rows.Next() {
		var a carbon
		if err := rows.Scan(&a.ID, &a.CarbonID, &a.SendID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.Duration); err != nil {
			return nil, err
		}
		files = append(files, a)
	}

	return files, nil
}

func (a *carbon) getCarbonFile(db *sql.DB) error {

	return db.QueryRow("SELECT id, carbon_id, send_id,  date, file_name, language, extension, size, sha1, duration FROM carbon WHERE carbon_id = $1",
		a.CarbonID).Scan(&a.ID, &a.CarbonID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1)
}

func (a *carbon) postCarbonFile(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO carbon(carbon_id, send_id,  date, file_name, language, extension, size, sha1, duration) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (carbon_id) DO UPDATE SET (carbon_id, send_id, date, file_name, language, extension, size, sha1, duration) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE carbon.carbon_id = $1 RETURNING id",
		a.CarbonID, a.SendID, a.Date, a.FileName, a.Language, a.Extension, a.Size, a.Sha1, a.Duration).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *carbon) deleteCarbonFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM carbon WHERE carbon_id=$1", a.CarbonID)

	return err
}