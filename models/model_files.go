package models

import (
	"database/sql"
	"encoding/json"
)

type Files struct {
	ID        int         `json:"id"`
	FileID    string      `json:"file_id"`
	Date      string      `json:"date"`
	Language  string      `json:"language"`
	FileName  string      `json:"file_name"`
	Extension string      `json:"extension"`
	Size      int64       `json:"size"`
	Sha1      string      `json:"sha1"`
	FileType  string      `json:"file_type"`
	MimeType  string      `json:"mime_type"`
	UID       string      `json:"uid"`
	WID       string      `json:"wid"`
	Props     interface{} `json:"properties"`
}

func FindFiles(db *sql.DB, key string, value string) ([]Files, error) {
	sqlStatement := `SELECT id, file_id, date, language, file_name, extension, size, sha1, file_type, mime_type, uid, wid, properties FROM files WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	o := []Files{}

	for rows.Next() {
		var a Files
		var properties []byte
		if err := rows.Scan(&a.ID, &a.FileID, &a.Date, &a.Language, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.MimeType, &a.UID, &a.WID, &properties); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &a.Props)
		o = append(o, a)
	}

	return o, nil
}

func GetFiles(db *sql.DB, start, count int) ([]Files, error) {
	rows, err := db.Query(
		"SELECT id, file_id, date, language, file_name, extension, size, sha1, file_type, mime_type, uid, wid, properties FROM files ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	o := []Files{}

	for rows.Next() {
		var a Files
		var properties []byte
		if err := rows.Scan(&a.ID, &a.FileID, &a.Date, &a.Language, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.MimeType, &a.UID, &a.WID, &properties); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &a.Props)
		o = append(o, a)
	}

	return o, nil
}

func (a *Files) GetFile(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT id, file_id, date, language, file_name, extension, size, sha1, file_type, mime_type, uid, wid, proertires FROM files WHERE file_id = $1",
		a.FileID).Scan(&a.ID, &a.FileID, &a.Date, &a.Language, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.MimeType, &a.UID, &a.WID, &properties)
	json.Unmarshal(properties, &a.Props)

	if err != nil {
		return err
	}

	return err
}

func (a *Files) PostFile(db *sql.DB) error {
	properties, _ := json.Marshal(a.Props)
	err := db.QueryRow(
		"INSERT INTO files(file_id, date, language, file_name, extension, size, sha1, file_type, mime_type, uid, wid, properties) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (sha1) DO UPDATE SET (file_id, date, language, file_name, extension, size, sha1, file_type, mime_type, uid, wid, properties) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) WHERE files.sha1 = $7 RETURNING id",
		&a.FileID, &a.Date, &a.Language, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.MimeType, &a.UID, &a.WID, &properties).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *Files) PostFileStatus(db *sql.DB, value, key string) error {
	_, err := db.Exec("UPDATE files SET properties = properties || json_build_object($3::text, $2::bool)::jsonb WHERE file_id=$1",
		a.FileID, value, key)

	return err
}

func (a *Files) PostFileValue(db *sql.DB, value string, key string) error {
	_, err := db.Exec("UPDATE files SET properties = properties || json_build_object($3::text, $2::text)::jsonb WHERE file_id=$1",
		a.FileID, value, key)

	return err
}

func (a *Files) PostFileJSON(db *sql.DB, value interface{}, key string) error {
	v, _ := json.Marshal(value)
	_, err := db.Exec("UPDATE files SET properties = properties || json_build_object($3::text, $2::jsonb)::jsonb WHERE file_id=$1",
		a.FileID, v, key)

	return err
}

func (a *Files) DeleteFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM files WHERE file_id=$1", a.FileID)

	return err
}