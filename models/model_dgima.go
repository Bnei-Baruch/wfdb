// model_dgima.go

package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type Dgims struct {
	ID        int                    `json:"id"`
	DgimaID   string                 `json:"dgima_id"`
	Date      string                 `json:"date"`
	FileName  string                 `json:"file_name"`
	Inpoints  []float32              `json:"inpoints"`
	Outpoints []float32              `json:"outpoints"`
	Parent    map[string]interface{} `json:"parent"`
	Line      map[string]interface{} `json:"line"`
	Original  map[string]interface{} `json:"original"`
	Proxy     map[string]interface{} `json:"proxy"`
	Wfstatus  map[string]interface{} `json:"wfstatus"`
}

func FindDgima(db *sql.DB, key string, value string) ([]Dgims, error) {
	sqlStatement := `SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY dgima_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	return objects, nil
}

func FindDgimaByJSON(db *sql.DB, ep string, key string, value string) ([]Dgims, error) {

	sqlStatement := fmt.Sprintf("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE %s ->> '%s' = '%s' ORDER BY dgima_id;", ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func FindDgimaBySHA1(db *sql.DB, value string) ([]Dgims, error) {

	sqlStatement := fmt.Sprintf("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE original->'format'->>'sha1' = '%s' ORDER BY dgima_id;", value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte
		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetDgima(db *sql.DB, start, count int) ([]Dgims, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima ORDER BY dgima_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	return objects, nil
}

func GetDgimaBySource(db *sql.DB, value string) ([]Dgims, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE wfstatus ->> 'removed' = 'false' AND parent ->> 'source' = $1 ORDER BY dgima_id", value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	return objects, nil
}

func GetFilesToDgima(db *sql.DB) ([]Dgims, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE wfstatus ->> 'removed' = 'false' AND parent ->> 'source' != 'cassette' ORDER BY dgima_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	return objects, nil
}

func GetCassetteFiles(db *sql.DB) ([]Dgims, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE wfstatus ->> 'removed' = 'false' AND parent ->> 'source' = 'cassette' ORDER BY dgima_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Dgims{}

	for rows.Next() {
		var d Dgims
		var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

		if err := rows.Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(inpoints, &d.Inpoints)
		json.Unmarshal(outpoints, &d.Outpoints)
		json.Unmarshal(parent, &d.Parent)
		json.Unmarshal(line, &d.Line)
		json.Unmarshal(original, &d.Original)
		json.Unmarshal(proxy, &d.Proxy)
		json.Unmarshal(wfstatus, &d.Wfstatus)
		objects = append(objects, d)
	}

	return objects, nil
}

func (d *Dgims) GetDgimaID(db *sql.DB) error {
	var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE dgima_id = $1",
		d.DgimaID).Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(inpoints, &d.Inpoints)
	json.Unmarshal(outpoints, &d.Outpoints)
	json.Unmarshal(parent, &d.Parent)
	json.Unmarshal(line, &d.Line)
	json.Unmarshal(original, &d.Original)
	json.Unmarshal(proxy, &d.Proxy)
	json.Unmarshal(wfstatus, &d.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (d *Dgims) GetDgimaByID(db *sql.DB) error {
	var inpoints, outpoints, parent, line, original, proxy, wfstatus []byte

	err := db.QueryRow("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE id = $1",
		d.ID).Scan(&d.ID, &d.DgimaID, &d.Date, &d.FileName, &inpoints, &outpoints, &parent, &line, &original, &proxy, &wfstatus)

	json.Unmarshal(inpoints, &d.Inpoints)
	json.Unmarshal(outpoints, &d.Outpoints)
	json.Unmarshal(parent, &d.Parent)
	json.Unmarshal(line, &d.Line)
	json.Unmarshal(original, &d.Original)
	json.Unmarshal(proxy, &d.Proxy)
	json.Unmarshal(wfstatus, &d.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (d *Dgims) PostDgimaID(db *sql.DB) error {
	parent, _ := json.Marshal(d.Parent)
	line, _ := json.Marshal(d.Line)
	original, _ := json.Marshal(d.Original)
	proxy, _ := json.Marshal(d.Proxy)
	wfstatus, _ := json.Marshal(d.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO dgima(dgima_id, date, file_name, inpoints, outpoints, parent, line, original, proxy, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (dgima_id) DO UPDATE SET (dgima_id, date, file_name, inpoints, outpoints, parent, line, original, proxy, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE dgima.dgima_id = $1 RETURNING id",
		d.DgimaID, d.Date, d.FileName, pq.Array(d.Inpoints), pq.Array(d.Outpoints), parent, line, original, proxy, wfstatus).Scan(&d.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *Dgims) PostDgimaJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE dgima SET ` + key + ` = $2 WHERE dgima_id=$1;`
	_, err := db.Exec(sqlStatement, d.DgimaID, v)

	return err
}

func (d *Dgims) PostDgimaValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE dgima SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE dgima_id=$1",
		d.DgimaID, value, key)

	return err
}

func (d *Dgims) DeleteDgimaID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dgima WHERE dgima_id=$1", d.DgimaID)

	return err
}
