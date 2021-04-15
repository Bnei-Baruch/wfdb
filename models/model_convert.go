package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Convert struct {
	ID        int                    `json:"id"`
	ConvertID string                 `json:"convert_id"`
	Name      string                 `json:"name"`
	Date      string                 `json:"date"`
	Progress  string                 `json:"progress"`
	State     string                 `json:"state"`
	Timestamp string                 `json:"timestamp"`
	Langcheck map[string]interface{} `json:"langcheck"`
}

func FindConvert(db *sql.DB, key string, value string) ([]Convert, error) {
	sqlStatement := `SELECT * FROM convert WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY convert_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Convert{}

	for rows.Next() {
		var i Convert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Langcheck)
		files = append(files, i)
	}

	return files, nil
}

func FindConvertByJSON(db *sql.DB, key string, value string) ([]Convert, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM convert WHERE line ->> '%s' = '%s' ORDER BY convert_id;", key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []Convert{}

	for rows.Next() {
		var i Convert
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

func GetConvert(db *sql.DB, start, count int) ([]Convert, error) {
	rows, err := db.Query(
		"SELECT * FROM convert ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Convert{}

	for rows.Next() {
		var i Convert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Langcheck)
		files = append(files, i)
	}

	return files, nil
}

func (i *Convert) GetConvertByID(db *sql.DB) error {
	var obj []byte

	err := db.QueryRow("SELECT * FROM convert WHERE convert_id = $1",
		i.ConvertID).Scan(&i.ID, &i.ConvertID, &i.Name, &i.Date, &i.Progress, &i.State, &i.Timestamp, &obj)

	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &i.Langcheck)

	return err
}

func (i *Convert) PostConvert(db *sql.DB) error {
	langcheck, _ := json.Marshal(i.Langcheck)

	err := db.QueryRow(
		"INSERT INTO convert(convert_id, name, date, progress, state, timestamp, langcheck) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (convert_id) DO UPDATE SET (convert_id, name, date, progress, state, timestamp, langcheck) = ($1, $2, $3, $4, $5, $6, $7) WHERE convert.convert_id = $1 RETURNING id",
		i.ConvertID, i.Name, i.Date, i.Progress, i.State, i.Timestamp, langcheck).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *Convert) PostConvertValue(db *sql.DB, key string, value string) error {

	sqlStatement := fmt.Sprintf("UPDATE convert SET %s='%s' WHERE convert_id='%s';", key, value, i.ConvertID)
	_, err := db.Exec(sqlStatement)

	return err
}

func (i *Convert) PostConvertJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE convert SET ` + key + ` = $2 WHERE convert_id=$1;`
	_, err := db.Exec(sqlStatement, i.ConvertID, v)

	return err
}

func (i *Convert) DeleteConvert(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM convert WHERE convert_id=$1", i.ConvertID)

	return err
}
