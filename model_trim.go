// model_trim.go

package main

import (
	"database/sql"
	"encoding/json"
)

type trim struct {
	ID    int						`json:"id"`
	TrimID  string  				`json:"trim_id"`
	Data map[string]interface{}		`json:"data"`
}

func findTrimes(db *sql.DB, key string, value string) ([]trim, error) {
	rows, err := db.Query(
		"SELECT id, trim_id, data FROM trim WHERE data @> json_build_object($1::text, $2::text)::jsonb",
		key, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trimes := []trim{}

	for rows.Next() {
		var t trim
		var obj []byte
		if err := rows.Scan(&t.ID, &t.TrimID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &t.Data)
		trimes = append(trimes, t)
	}

	return trimes, nil
}

func getTrimes(db *sql.DB, start, count int) ([]trim, error) {
	rows, err := db.Query(
		"SELECT id, trim_id, data FROM trim ORDER BY trim_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trimes := []trim{}

	for rows.Next() {
		var t trim
		var obj []byte
		if err := rows.Scan(&t.ID, &t.TrimID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &t.Data)
		trimes = append(trimes, t)
	}

	return trimes, nil
}

func (t *trim) getTrim(db *sql.DB) error {
	var obj []byte
	err := db.QueryRow("SELECT data FROM trim WHERE trim_id = $1",
		t.TrimID).Scan(&obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &t.Data)

	return err
}

func (t *trim) postTrim(db *sql.DB) error {
	v, _ := json.Marshal(t.Data)

	err := db.QueryRow(
		"INSERT INTO trim(trim_id, data) VALUES($1, $2) ON CONFLICT (trim_id) DO UPDATE SET (data) = ($2) WHERE trim.trim_id = $1 RETURNING id",
		t.TrimID, v).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *trim) updateTrim(db *sql.DB) error {
	v, _ := json.Marshal(t.Data)
	_, err :=
		db.Exec("UPDATE trim SET data=$2 WHERE trim_id=$1",
			t.TrimID, v)

	return err
}

func (t *trim) deleteTrim(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM trim WHERE trim_id=$1", t.TrimID)

	return err
}