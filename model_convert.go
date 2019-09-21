// model_convert.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type convert struct {
	ID    int							`json:"id"`
	ConvertID  string 				 	`json:"convert_id"`
	Name string							`json:"name"`
	Date string							`json:"date"`
	Progress string						`json:"progress"`
	State string						`json:"state"`
	Timestamp string					`json:"timestamp"`
	Langcheck map[string]interface{}	`json:"langcheck"`
}

func findConvert(db *sql.DB, key string, value string) ([]convert, error) {
	sqlStatement := `SELECT * FROM convert WHERE `+key+` LIKE '%`+value+`%' ORDER BY convert_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []convert{}

	for rows.Next() {
		var i convert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Langcheck)
		files = append(files, i)
	}

	return files, nil
}

func findConvertByJSON(db *sql.DB, key string, value string) ([]convert, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM convert WHERE line ->> '%s' = '%s' ORDER BY convert_id;", key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []convert{}

	for rows.Next() {
		var i convert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Langcheck)
		files = append(files, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func getConvert(db *sql.DB, start, count int) ([]convert, error) {
	rows, err := db.Query(
		"SELECT * FROM convert ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []convert{}

	for rows.Next() {
		var i convert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Langcheck)
		files = append(files, i)
	}

	return files, nil
}

func (i *convert) getConvertByID(db *sql.DB) error {
	var obj []byte

	err := db.QueryRow("SELECT * FROM convert WHERE convert_id = $1",
		i.ConvertID).Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj)

	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &i.Langcheck)

	return err
}

func (i *convert) postConvert(db *sql.DB) error {
	langcheck, _ := json.Marshal(i.Langcheck)

	err := db.QueryRow(
		"INSERT INTO convert(convert_id, name, date, progress, state, timestamp, langcheck) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (convert_id) DO UPDATE SET (convert_id, name, date, progress, state, timestamp, langcheck) = ($1, $2, $3, $4, $5, $6, $7) WHERE convert.convert_id = $1 RETURNING id",
		i.ID, i.ConvertID, i.Name, i.Date, i.Progress, i.State, i.Timestamp, langcheck).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *convert) deleteConvert(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM convert WHERE convert_id=$1", i.ConvertID)

	return err
}