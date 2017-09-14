// model_capture.go

package main

import (
	"database/sql"
	"encoding/json"
)

type capture struct {
	ID    int						`json:"id"`
	CaptureID  string  				`json:"capture_id"`
	Data map[string]interface{}		`json:"data"`
}

func findCaptures(db *sql.DB, key string, value string) ([]capture, error) {
	rows, err := db.Query(
		"SELECT id, capture_id, data FROM capture WHERE data @> json_build_object($1::text, $2::text)::jsonb",
		key, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	captures := []capture{}

	for rows.Next() {
		var c capture
		var obj []byte
		if err := rows.Scan(&c.ID, &c.CaptureID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &c.Data)
		captures = append(captures, c)
	}

	return captures, nil
}

func getCaptures(db *sql.DB, start, count int) ([]capture, error) {
	rows, err := db.Query(
		"SELECT id, capture_id, data FROM capture ORDER BY capture_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	captures := []capture{}

	for rows.Next() {
		var c capture
		var obj []byte
		if err := rows.Scan(&c.ID, &c.CaptureID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &c.Data)
		captures = append(captures, c)
	}

	return captures, nil
}

func (c *capture) getCapture(db *sql.DB) error {
	var obj []byte
	err := db.QueryRow("SELECT data FROM capture WHERE capture_id = $1",
		c.CaptureID).Scan(&obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &c.Data)

	return err
}

func (c *capture) postCapture(db *sql.DB) error {
	v, _ := json.Marshal(c.Data)

	err := db.QueryRow(
		"INSERT INTO capture(capture_id, data) VALUES($1, $2) ON CONFLICT (capture_id) DO UPDATE SET (data) = ($2) WHERE capture.capture_id = $1 RETURNING id",
		c.CaptureID, v).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *capture) updateCapture(db *sql.DB) error {
	v, _ := json.Marshal(c.Data)
	_, err :=
		db.Exec("UPDATE capture SET data=$2 WHERE capture_id=$1",
			c.CaptureID, v)

	return err
}

func (c *capture) deleteCapture(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM capture WHERE capture_id=$1", c.CaptureID)

	return err
}