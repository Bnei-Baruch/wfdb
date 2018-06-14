// model_ingest.go

package main

import (
	"database/sql"
	"encoding/json"
)

type ingest struct {
	ID    int			`json:"id"`
	CaptureID  string  	`json:"capture_id"`
	Capsrc  string  	`json:"capture_src"`
	Date string			`json:"date"`
	StartName string	`json:"start_name"`
	StopName string		`json:"stop_name"`
	Sha1 string			`json:"sha1"`
	Line map[string]interface{}			`json:"line"`
	Original map[string]interface{}		`json:"original"`
	Proxy map[string]interface{}		`json:"proxy"`
	Wfstatus map[string]interface{}		`json:"wfstatus"`
}

func findIngest(db *sql.DB, key string, value string) ([]ingest, error) {
	sqlStatement := `SELECT * FROM ingest WHERE `+key+` LIKE '%`+value+`%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []ingest{}

	for rows.Next() {
		var i ingest
		var line, original, proxy, wfstatus []byte
		if err := rows.Scan(&i.ID, &i.CaptureID, &i.Capsrc, &i.Date, &i.StartName, &i.StopName, &i.Sha1, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &i.Line)
		json.Unmarshal(original, &i.Original)
		json.Unmarshal(proxy, &i.Proxy)
		json.Unmarshal(wfstatus, &i.Wfstatus)
		objects = append(objects, i)
	}

	return objects, nil
}

func getIngest(db *sql.DB, start, count int) ([]ingest, error) {
	rows, err := db.Query(
		"SELECT * FROM ingest ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []ingest{}

	for rows.Next() {
		var i ingest
		var line, original, proxy, wfstatus []byte
		if err := rows.Scan(&i.ID, &i.CaptureID, &i.Capsrc, &i.Date, &i.StartName, &i.StopName, &i.Sha1, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &i.Line)
		json.Unmarshal(original, &i.Original)
		json.Unmarshal(proxy, &i.Proxy)
		json.Unmarshal(wfstatus, &i.Wfstatus)
		objects = append(objects, i)
	}

	return objects, nil
}

func (i *ingest) getIngestID(db *sql.DB) error {
	var line []byte
	var original []byte
	var proxy []byte
	var wfstatus []byte

	err := db.QueryRow("SELECT * FROM ingest WHERE capture_id = $1",
		i.CaptureID).Scan(&i.ID, &i.CaptureID, &i.Capsrc, &i.Date, &i.StartName, &i.StopName, &i.Sha1, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(line, &i.Line)
	json.Unmarshal(original, &i.Original)
	json.Unmarshal(proxy, &i.Proxy)
	json.Unmarshal(wfstatus, &i.Wfstatus)

	if err != nil {
		return err
	}

	return err
}

func (i *ingest) postIngestID(db *sql.DB) error {
	line, _ := json.Marshal(i.Line)
	original, _ := json.Marshal(i.Original)
	proxy, _ := json.Marshal(i.Proxy)
	wfstatus, _ := json.Marshal(i.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO ingest(capture_id, capture_src, date, start_name, stop_name, sha1, line, original, proxy, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (capture_id) DO UPDATE SET (capture_id, capture_src, date, start_name, stop_name, sha1, line, original, proxy, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE ingest.capture_id = $1 RETURNING id",
		i.CaptureID, i.Capsrc, i.Date, i.StartName, i.StopName, i.Sha1, line, original, proxy, wfstatus).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *ingest) postIngestJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE ingest SET `+key+` = $2 WHERE capture_id=$1;`
	//sqlStatement := `UPDATE ingest SET wfstatus = wtstatus || '{"`+key+`": $2}' WHERE capture_id=$1;`
	_, err := db.Exec(sqlStatement, i.CaptureID, v)

	return err
}

func (i *ingest) postIngestValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE ingest SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE capture_id=$1",
		i.CaptureID, value, key)

	return err
}

func (i *ingest) deleteIngestID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM ingest WHERE capture_id=$1", i.CaptureID)

	return err
}