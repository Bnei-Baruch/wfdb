package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Clouds struct {
	ID        int         `json:"id"`
	OID       string      `json:"oid"`
	Date      string      `json:"date"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Extension string      `json:"extension"`
	Language  string      `json:"language"`
	Source    string      `json:"source"`
	UID       string      `json:"uid"`
	WID       string      `json:"wid"`
	Pattern   string      `json:"pattern"`
	Props     interface{} `json:"properties"`
	Url       string      `json:"url"`
}

func FindCloud(db *sql.DB, values url.Values) ([]Clouds, error) {

	var where []string
	sqlStatement := `SELECT * FROM cloud WHERE properties['removed'] = 'false'`
	limit := "10"
	offset := "0"

	for k, v := range values {
		if k == "limit" {
			limit = v[0]
			continue
		}
		if k == "offset" {
			offset = v[0]
			continue
		}
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, k, v[0]))
	}

	if len(where) > 0 {
		sqlStatement = sqlStatement + ` AND ` + strings.Join(where, " AND ")
	}

	sqlStatement = sqlStatement + fmt.Sprintf(` ORDER BY oid DESC LIMIT %s OFFSET %s`, limit, offset)

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Clouds{}

	for rows.Next() {
		var t Clouds
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindCloudByJSON(db *sql.DB, ep string, key string, value string) ([]Clouds, error) {

	sqlStatement := fmt.Sprintf(`SELECT * FROM cloud WHERE %s['%s'] = '"%s"' ORDER BY oid;`, ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Clouds{}

	for rows.Next() {
		var t Clouds
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetListClouds(db *sql.DB, start, count int) ([]Clouds, error) {
	rows, err := db.Query(
		"SELECT * FROM cloud ORDER BY oid DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Clouds{}

	for rows.Next() {
		var t Clouds
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Clouds) GetCloudID(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT * FROM cloud WHERE oid = $1",
		t.OID).Scan(
		&t.ID, &t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Clouds) GetCloudByID(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT * FROM cloud WHERE id = $1",
		t.ID).Scan(&t.ID, &t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Clouds) PostCloudID(db *sql.DB) error {
	properties, _ := json.Marshal(t.Props)

	err := db.QueryRow(
		"INSERT INTO cloud(oid, date, name, type, extension, language, source, uid, wid, pattern, properties, url) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (oid) DO UPDATE SET (oid, date, name, type, extension, language, source, uid, wid, pattern, properties, url) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) WHERE cloud.oid = $1 RETURNING id",
		&t.OID, &t.Date, &t.Name, &t.Type, &t.Extension, &t.Language, &t.Source, &t.UID, &t.WID, &t.Pattern, &properties, &t.Url).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *Clouds) PostCloudJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE cloud SET ` + key + ` = $2 WHERE oid=$1;`
	_, err := db.Exec(sqlStatement, t.OID, v)

	return err
}

func (t *Clouds) SetCloudJSON(db *sql.DB, value interface{}, key string, prop string) error {

	v, _ := json.Marshal(value)
	sqlCmd := "UPDATE cloud SET " + prop + " = " + prop + " || json_build_object($3::text, $2::jsonb)::jsonb WHERE oid=$1"
	_, err := db.Exec(sqlCmd, t.OID, v, key)

	return err
}

func (t *Clouds) PostCloudStatus(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE cloud SET properties = properties || json_build_object($3::text, $2::bool)::jsonb WHERE oid=$1",
		t.OID, value, key)

	return err
}

func (t *Clouds) PostCloudProp(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE cloud SET properties = properties || json_build_object($3::text, $2::text)::jsonb WHERE oid=$1",
		t.OID, value, key)

	return err
}

func (t *Clouds) DeleteCloudID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM cloud WHERE oid=$1", t.OID)

	return err
}
