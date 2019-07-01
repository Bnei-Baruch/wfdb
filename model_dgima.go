// model_dgima.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type dgima struct {
	ID    int			`json:"id"`
	DgimaID  string  	`json:"dgima_id"`
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

func findDgima(db *sql.DB, key string, value string) ([]dgima, error) {
	sqlStatement := `SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE `+key+` LIKE '%`+value+`%' ORDER BY dgima_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func findDgimaByJSON(db *sql.DB, ep string, key string, value string) ([]dgima, error) {

	sqlStatement := fmt.Sprintf("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE %s ->> '%s' = '%s' ORDER BY dgima_id;", ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func findDgimaBySHA1(db *sql.DB, value string) ([]dgima, error) {

	sqlStatement := fmt.Sprintf("SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE original->'format'->>'sha1' = '%s' ORDER BY dgima_id;", value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func getDgima(db *sql.DB, start, count int) ([]dgima, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima ORDER BY dgima_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func getDgimaBySource(db *sql.DB, value string) ([]dgima, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE wfstatus ->> 'removed' = 'false' AND parent ->> 'source' = $1 ORDER BY dgima_id", value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func getFilesToDgima(db *sql.DB) ([]dgima, error) {
	rows, err := db.Query(
		"SELECT id, dgima_id, date, file_name, array_to_json(inpoints), array_to_json(outpoints), parent, line, original, proxy, wfstatus FROM dgima WHERE wfstatus ->> 'removed' = 'false' ORDER BY dgima_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []dgima{}

	for rows.Next() {
		var d dgima
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

func (d *dgima) getDgimaID(db *sql.DB) error {
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

func (d *dgima) getDgimaByID(db *sql.DB) error {
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


func (d *dgima) postDgimaID(db *sql.DB) error {
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

func (d *dgima) postDgimaJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE dgima SET `+key+` = $2 WHERE dgima_id=$1;`
	_, err := db.Exec(sqlStatement, d.DgimaID, v)

	return err
}

func (d *dgima) postDgimaValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE dgima SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE dgima_id=$1",
		d.DgimaID, value, key)

	return err
}

func (d *dgima) deleteDgimaID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dgima WHERE dgima_id=$1", d.DgimaID)

	return err
}