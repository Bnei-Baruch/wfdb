// model_capture.go

package main

import (
	"database/sql"
	"encoding/json"
)

type capture struct {
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

func findCapture(db *sql.DB, key string, value string) ([]capture, error) {
	sqlStatement := `SELECT * FROM capture WHERE `+key+` LIKE '%`+value+`%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []capture{}

	for rows.Next() {
		var c capture
		var line, original, proxy, wfstatus []byte
		if err := rows.Scan(&c.ID, &c.CaptureID, &c.Capsrc, &c.Date, &c.StartName, &c.StopName, &c.Sha1, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &c.Line)
		json.Unmarshal(original, &c.Original)
		json.Unmarshal(proxy, &c.Proxy)
		json.Unmarshal(wfstatus, &c.Wfstatus)
		objects = append(objects, c)
	}

	return objects, nil
}

func getCapture(db *sql.DB, start, count int) ([]capture, error) {
	rows, err := db.Query(
		"SELECT * FROM capture ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []capture{}

	for rows.Next() {
		var c capture
		var line, original, proxy, wfstatus []byte
		if err := rows.Scan(&c.ID, &c.CaptureID, &c.Capsrc, &c.Date, &c.StartName, &c.StopName, &c.Sha1, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &c.Line)
		json.Unmarshal(original, &c.Original)
		json.Unmarshal(proxy, &c.Proxy)
		json.Unmarshal(wfstatus, &c.Wfstatus)
		objects = append(objects, c)
	}

	return objects, nil
}

func (c *capture) getCaptureID(db *sql.DB) error {
	var line []byte
	var original []byte
	var proxy []byte
	var wfstatus []byte

	err := db.QueryRow("SELECT * FROM capture WHERE capture_id = $1",
		c.CaptureID).Scan(&c.ID, &c.CaptureID, &c.Capsrc, &c.Date, &c.StartName, &c.StopName, &c.Sha1, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(line, &c.Line)
	json.Unmarshal(original, &c.Original)
	json.Unmarshal(proxy, &c.Proxy)
	json.Unmarshal(wfstatus, &c.Wfstatus)

	if err != nil {
		return err
	}

	return err
}

func (c *capture) postCaptureID(db *sql.DB) error {
	line, _ := json.Marshal(c.Line)
	original, _ := json.Marshal(c.Original)
	proxy, _ := json.Marshal(c.Proxy)
	wfstatus, _ := json.Marshal(c.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO capture(capture_id, capture_src, date, start_name, stop_name, sha1, line, original, proxy, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (capture_id) DO UPDATE SET (capture_id, capture_src, date, start_name, stop_name, sha1, line, original, proxy, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE capture.capture_id = $1 RETURNING id",
		c.CaptureID, c.Capsrc, c.Date, c.StartName, c.StopName, c.Sha1, line, original, proxy, wfstatus).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *capture) postCaptureJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE capture SET `+key+` = $2 WHERE capture_id=$1;`
	//sqlStatement := `UPDATE capture SET wfstatus = wtstatus || '{"`+key+`": $2}' WHERE capture_id=$1;`
	_, err := db.Exec(sqlStatement, c.CaptureID, v)

	return err
}

func (c *capture) postCaptureValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE capture SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE capture_id=$1",
		c.CaptureID, value, key)

	return err
}

func (c *capture) deleteCaptureID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM capture WHERE capture_id=$1", c.CaptureID)

	return err
}