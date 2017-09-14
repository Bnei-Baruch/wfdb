// model_trimmer.go

package main

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
)

type trimmer struct {
	ID    int			`json:"id"`
	TrimID  string  	`json:"trim_id"`
	Date string			`json:"date"`
	FileName string		`json:"file_name"`
	Inpoints []float32	`json:"inpoints"`
	Outpoints []float32	`json:"outpoints"`
	Parent map[string]interface{}		`json:"parent"`
	Line map[string]interface{}			`json:"line"`
	Original map[string]interface{}		`json:"original"`
	Proxy map[string]interface{}		`json:"proxy"`
	Wfstatus map[string]interface{}		`json:"wfstatus"`
}

func findTrimmer(db *sql.DB, key string, value string) ([]trimmer, error) {
	sqlStatement := `SELECT id, trim_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM trimmer WHERE `+key+` LIKE '%`+value+`' ORDER BY trim_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []trimmer{}

	for rows.Next() {
		var t trimmer
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.TrimID, &t.Date, &t.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &t.Inpoints)
		json.Unmarshal(outpoints, &t.Outpoints)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func getTrimmer(db *sql.DB, start, count int) ([]trimmer, error) {
	rows, err := db.Query(
		"SELECT id, trim_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM trimmer ORDER BY trim_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []trimmer{}

	for rows.Next() {
		var t trimmer
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&t.ID, &t.TrimID, &t.Date, &t.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &t.Inpoints)
		json.Unmarshal(outpoints, &t.Outpoints)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func getFilesToTrim(db *sql.DB) ([]trimmer, error) {
	rows, err := db.Query(
		"SELECT id, trim_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM trimmer WHERE wfstatus ->> 'removed' = 'false' ORDER BY trim_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []trimmer{}

	for rows.Next() {
		var t trimmer
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&t.ID, &t.TrimID, &t.Date, &t.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &t.Inpoints)
		json.Unmarshal(outpoints, &t.Outpoints)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *trimmer) getTrimmerID(db *sql.DB) error {
	var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, trim_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM trimmer WHERE trim_id = $1",
		t.TrimID).Scan(&t.ID, &t.TrimID, &t.Date, &t.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(inpoints, &t.Inpoints)
	json.Unmarshal(outpoints, &t.Outpoints)
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

func (t *trimmer) getTrimmerByID(db *sql.DB) error {
	var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, trim_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM trimmer WHERE id = $1",
		t.ID).Scan(&t.ID, &t.TrimID, &t.Date, &t.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(inpoints, &t.Inpoints)
	json.Unmarshal(outpoints, &t.Outpoints)
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


func (t *trimmer) postTrimmerID(db *sql.DB) error {
	parent, _ := json.Marshal(t.Parent)
	line, _ := json.Marshal(t.Line)
	original, _ := json.Marshal(t.Original)
	proxy, _ := json.Marshal(t.Proxy)
	wfstatus, _ := json.Marshal(t.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO trimmer(trim_id, date, file_name, inpoints, outpoints, parent, line, original, proxy, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (trim_id) DO UPDATE SET (trim_id, date, file_name, inpoints, outpoints, parent, line, original, proxy, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE trimmer.trim_id = $1 RETURNING id",
		t.TrimID, t.Date, t.FileName, pq.Array(t.Inpoints), pq.Array(t.Outpoints), parent, line, original, proxy, wfstatus).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *trimmer) postTrimmerJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE trimmer SET `+key+` = $2 WHERE trim_id=$1;`
	_, err := db.Exec(sqlStatement, t.TrimID, v)

	return err
}

func (t *trimmer) postTrimmerValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE trimmer SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE trim_id=$1",
		t.TrimID, value, key)

	return err
}

func (t *trimmer) deleteTrimmerID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM trimmer WHERE trim_id=$1", t.TrimID)

	return err
}