// model_labels.go

package main

import (
	"database/sql"
)

type labels struct {
	ID    int			`json:"id"`
	Date string			`json:"date"`
	Lecturer string		`json:"lecturer"`
	Subject string		`json:"subject"`
	Language string		`json:"language"`
	Location string		`json:"location"`
	ContentType string	`json:"content_type"`
	Cassete_type string	`json:"cassete_type"`
	Mof string			`json:"mof"`
	Duration string		`json:"duration"`
	Archive_place string	`json:"archive_place"`
	Comments string		`json:"comments"`
	Bar_code  string  	`json:"bar_code"`
}

func findLabels(db *sql.DB, key string, value string) ([]labels, error) {
	sqlStatement := `SELECT * FROM labels WHERE `+key+` LIKE '%`+value+`' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []labels{}

	for rows.Next() {
		var l labels
		if err := rows.Scan(&l.ID, &l.Date, &l.Lecturer, &l.Subject, &l.Language, &l.Location, &l.ContentType, &l.Cassete_type, &l.Mof, &l.Duration, &l.Archive_place, &l.Comments, &l.Bar_code); err != nil {
			return nil, err
		}
		objects = append(objects, l)
	}

	return objects, nil
}

func getLabels(db *sql.DB, start, count int) ([]labels, error) {
	rows, err := db.Query(
		"SELECT * FROM labels ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []labels{}

	for rows.Next() {
		var l labels
		if err := rows.Scan(&l.ID, &l.Date, &l.Lecturer, &l.Subject, &l.Language, &l.Location, &l.ContentType, &l.Cassete_type, &l.Mof, &l.Duration, &l.Archive_place, &l.Comments, &l.Bar_code); err != nil {
			return nil, err
		}
		objects = append(objects, l)
	}

	return objects, nil
}

func (l *labels) getLabel(db *sql.DB) error {

	return db.QueryRow("SELECT * FROM labels WHERE id = $1",
		l.ID).Scan(&l.ID, &l.Date, &l.Lecturer, &l.Subject, &l.Language, &l.Location, &l.ContentType, &l.Cassete_type, &l.Mof, &l.Duration, &l.Archive_place, &l.Comments, &l.Bar_code)
}