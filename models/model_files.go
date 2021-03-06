package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Files struct {
	ID        int         `json:"id"`
	FileID    string      `json:"file_id"`
	Date      string      `json:"date"`
	FileName  string      `json:"file_name"`
	Extension string      `json:"extension"`
	Size      int64       `json:"size"`
	Sha1      string      `json:"sha1"`
	FileType  string      `json:"file_type"`
	Language  string      `json:"language"`
	MimeType  string      `json:"mime_type"`
	UID       string      `json:"uid"`
	WID       string      `json:"wid"`
	Props     interface{} `json:"properties"`
	ProductID string      `json:"product_id"`
}

func FindFiles(db *sql.DB, values url.Values) ([]Files, error) {

	var where []string

	for k, v := range values {
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, k, v[0]))
	}

	sqlStatement := ("SELECT * FROM files WHERE properties ->> 'removed' = 'false' AND " + strings.Join(where, " AND "))
	//sqlStatement := ("SELECT * FROM products WHERE properties ->> 'removed' = 'false' AND " + strings.Join(where, " AND "))
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	o := []Files{}

	for rows.Next() {
		var a Files
		var properties []byte
		if err := rows.Scan(&a.ID, &a.FileID, &a.Date, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.Language, &a.MimeType, &a.UID, &a.WID, &properties, &a.ProductID); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &a.Props)
		o = append(o, a)
	}

	return o, nil
}

func GetFiles(db *sql.DB, start, count int) ([]Files, error) {
	rows, err := db.Query(
		"SELECT * FROM files ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	o := []Files{}

	for rows.Next() {
		var a Files
		var properties []byte
		if err := rows.Scan(&a.ID, &a.FileID, &a.Date, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.Language, &a.MimeType, &a.UID, &a.WID, &properties, &a.ProductID); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &a.Props)
		o = append(o, a)
	}

	return o, nil
}

func GetActiveFiles(db *sql.DB, language string, product_id string) ([]Files, error) {
	rows, err := db.Query(
		"SELECT * FROM files WHERE language = $1 AND product_id = $2 ORDER BY file_id", language, product_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Files{}

	for rows.Next() {
		var a Files
		var properties []byte
		if err := rows.Scan(
			&a.ID, &a.FileID, &a.Date, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.Language, &a.MimeType, &a.UID, &a.WID, &properties, &a.ProductID); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &a.Props)
		objects = append(objects, a)
	}

	return objects, nil
}

func (a *Files) GetFile(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT * FROM files WHERE file_id = $1",
		a.FileID).Scan(&a.ID, &a.FileID, &a.Date, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.Language, &a.MimeType, &a.UID, &a.WID, &properties, &a.ProductID)
	json.Unmarshal(properties, &a.Props)

	if err != nil {
		return err
	}

	return err
}

func (a *Files) PostFile(db *sql.DB) error {
	properties, _ := json.Marshal(a.Props)
	err := db.QueryRow(
		"INSERT INTO files(file_id, date, file_name, extension, size, sha1, file_type, language,  mime_type, uid, wid, properties, product_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT (sha1) DO UPDATE SET (file_id, date, file_name, extension, size, sha1, file_type, language,  mime_type, uid, wid, properties, product_id) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) WHERE files.sha1 = $7 RETURNING id",
		&a.FileID, &a.Date, &a.FileName, &a.Extension, &a.Size, &a.Sha1, &a.FileType, &a.Language, &a.MimeType, &a.UID, &a.WID, &properties, &a.ProductID).Scan(&a.ID)

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
