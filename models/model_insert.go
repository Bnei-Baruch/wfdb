package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Insert struct {
	ID         int                    `json:"id"`
	InsertID   string                 `json:"insert_id"`
	InsertName string                 `json:"insert_name"`
	Date       string                 `json:"date"`
	FileName   string                 `json:"file_name"`
	Extension  string                 `json:"extension"`
	Size       int64                  `json:"size"`
	Sha1       string                 `json:"sha1"`
	Language   string                 `json:"language"`
	InsertType string                 `json:"insert_type"`
	SendID     string                 `json:"send_id"`
	UploadType string                 `json:"upload_type"`
	Line       map[string]interface{} `json:"line"`
}

func FindInsertFiles(db *sql.DB, key string, value string) ([]Insert, error) {
	sqlStatement := `SELECT * FROM insert WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY insert_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Insert{}

	for rows.Next() {
		var i Insert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.InsertID, &i.InsertName, &i.Date, &i.FileName, &i.Extension, &i.Size, &i.Sha1, &i.Language, &i.InsertType, &i.SendID, &i.UploadType, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Line)
		files = append(files, i)
	}

	return files, nil
}

func FindInsertByJSON(db *sql.DB, key string, value string) ([]Insert, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM insert WHERE line ->> '%s' = '%s' ORDER BY insert_id;", key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []Insert{}

	for rows.Next() {
		var i Insert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.InsertID, &i.InsertName, &i.Date, &i.FileName, &i.Extension, &i.Size, &i.Sha1, &i.Language, &i.InsertType, &i.SendID, &i.UploadType, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Line)
		files = append(files, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func GetInsertFiles(db *sql.DB, start, count int) ([]Insert, error) {
	rows, err := db.Query(
		"SELECT * FROM insert ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []Insert{}

	for rows.Next() {
		var i Insert
		var obj []byte
		if err := rows.Scan(&i.ID, &i.InsertID, &i.InsertName, &i.Date, &i.FileName, &i.Extension, &i.Size, &i.Sha1, &i.Language, &i.InsertType, &i.SendID, &i.UploadType, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &i.Line)
		files = append(files, i)
	}

	return files, nil
}

func (i *Insert) GetInsertFile(db *sql.DB) error {
	var obj []byte

	err := db.QueryRow("SELECT * FROM insert WHERE insert_id = $1",
		i.InsertID).Scan(&i.ID, &i.InsertID, &i.InsertName, &i.Date, &i.FileName, &i.Extension, &i.Size, &i.Sha1, &i.Language, &i.InsertType, &i.SendID, &i.UploadType, &obj)

	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &i.Line)

	return err
}

func (i *Insert) PostInsertFile(db *sql.DB) error {
	line, _ := json.Marshal(i.Line)

	err := db.QueryRow(
		"INSERT INTO insert(insert_id, insert_name, date, file_name, extension, size, sha1, language, insert_type, send_id, upload_type, line) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (insert_id) DO UPDATE SET (insert_id, insert_name, date, file_name, extension, size, sha1, language, insert_type, send_id, upload_type, line) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) WHERE insert.insert_id = $1 RETURNING id",
		i.InsertID, i.InsertName, i.Date, i.FileName, i.Extension, i.Size, i.Sha1, i.Language, i.InsertType, i.SendID, i.UploadType, line).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *Insert) DeleteInsertFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM insert WHERE insert_id=$1", i.InsertID)

	return err
}
