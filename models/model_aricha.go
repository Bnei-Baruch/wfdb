package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Aricha struct {
	ID       int                    `json:"id"`
	ArichaID string                 `json:"aricha_id"`
	Date     string                 `json:"date"`
	FileName string                 `json:"file_name"`
	Parent   map[string]interface{} `json:"parent"`
	Line     map[string]interface{} `json:"line"`
	Original map[string]interface{} `json:"original"`
	Proxy    map[string]interface{} `json:"proxy"`
	Wfstatus map[string]interface{} `json:"wfstatus"`
}

func FindAricha(db *sql.DB, key string, value string) ([]Aricha, error) {
	sqlStatement := `SELECT * FROM aricha WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY aricha_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Aricha{}

	for rows.Next() {
		var t Aricha
		var parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindArichaByJSON(db *sql.DB, ep string, key string, value string) ([]Aricha, error) {

	sqlStatement := fmt.Sprintf("SELECT id, aricha_id, date, file_name, parent, line, original, proxy, wfstatus FROM aricha WHERE %s ->> '%s' = '%s' ORDER BY aricha_id;", ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Aricha{}

	for rows.Next() {
		var t Aricha
		var parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func FindArichaBySHA1(db *sql.DB, value string) ([]Aricha, error) {

	sqlStatement := fmt.Sprintf("SELECT id, aricha_id, date, file_name, parent, line, original, proxy, wfstatus FROM aricha WHERE original->'format'->>'sha1' = '%s' ORDER BY aricha_id;", value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Aricha{}

	for rows.Next() {
		var t Aricha
		var parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetAricha(db *sql.DB, start, count int) ([]Aricha, error) {
	rows, err := db.Query(
		"SELECT * FROM aricha ORDER BY aricha_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Aricha{}

	for rows.Next() {
		var t Aricha
		var parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func GetBdika(db *sql.DB) ([]Aricha, error) {
	rows, err := db.Query(
		"SELECT * FROM aricha WHERE wfstatus ->> 'removed' = 'false' ORDER BY aricha_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Aricha{}

	for rows.Next() {
		var t Aricha
		var parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Aricha) GetArichaID(db *sql.DB) error {
	var parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, aricha_id, date, file_name, parent, line, original, proxy, wfstatus FROM aricha WHERE aricha_id = $1",
		t.ArichaID).Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(original, &t.Original)
	json.Unmarshal(proxy, &t.Proxy)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Aricha) GetArichaByID(db *sql.DB) error {
	var parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, aricha_id, date, file_name, parent, line, original, proxy, wfstatus FROM aricha WHERE id = $1",
		t.ID).Scan(&t.ID, &t.ArichaID, &t.Date, &t.FileName, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(original, &t.Original)
	json.Unmarshal(proxy, &t.Proxy)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Aricha) PostArichaID(db *sql.DB) error {
	parent, _ := json.Marshal(t.Parent)
	line, _ := json.Marshal(t.Line)
	original, _ := json.Marshal(t.Original)
	proxy, _ := json.Marshal(t.Proxy)
	wfstatus, _ := json.Marshal(t.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO aricha(aricha_id, date, file_name, parent, line, original, proxy, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (aricha_id) DO UPDATE SET (aricha_id, date, file_name, parent, line, original, proxy, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE aricha.aricha_id = $1 RETURNING id",
		t.ArichaID, t.Date, t.FileName, parent, line, original, proxy, wfstatus).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *Aricha) PostArichaJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE aricha SET ` + key + ` = $2 WHERE aricha_id=$1;`
	_, err := db.Exec(sqlStatement, t.ArichaID, v)

	return err
}

func (t *Aricha) PostArichaValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE aricha SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE aricha_id=$1",
		t.ArichaID, value, key)

	return err
}

func (t *Aricha) DeleteArichaID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM aricha WHERE aricha_id=$1", t.ArichaID)

	return err
}
