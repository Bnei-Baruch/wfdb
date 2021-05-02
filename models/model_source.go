package models

import (
	"database/sql"
	"encoding/json"
)

type Source struct {
	ID       int                    `json:"id"`
	SourceID string                 `json:"source_id"`
	Date     string                 `json:"date"`
	FileName string                 `json:"file_name"`
	Sha1     string                 `json:"sha1"`
	Line     map[string]interface{} `json:"line"`
	Source   map[string]interface{} `json:"source"`
	Wfstatus map[string]interface{} `json:"wfstatus"`
}

func FindSource(db *sql.DB, key string, value string) ([]Source, error) {
	sqlStatement := `SELECT * FROM source WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Source{}

	for rows.Next() {
		var i Source
		var line, source, wfstatus []byte
		if err := rows.Scan(&i.ID, &i.SourceID, &i.Date, &i.FileName, &i.Sha1, &line, &source, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &i.Line)
		json.Unmarshal(source, &i.Source)
		json.Unmarshal(wfstatus, &i.Wfstatus)
		objects = append(objects, i)
	}

	return objects, nil
}

func GetSource(db *sql.DB, start, count int) ([]Source, error) {
	rows, err := db.Query(
		"SELECT * FROM source ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Source{}

	for rows.Next() {
		var i Source
		var line, source, wfstatus []byte
		if err := rows.Scan(&i.ID, &i.SourceID, &i.Date, &i.FileName, &i.Sha1, &line, &source, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(line, &i.Line)
		json.Unmarshal(source, &i.Source)
		json.Unmarshal(wfstatus, &i.Wfstatus)
		objects = append(objects, i)
	}

	return objects, nil
}

func (i *Source) GetSourceID(db *sql.DB) error {
	var line []byte
	var source []byte
	var wfstatus []byte

	err := db.QueryRow("SELECT * FROM source WHERE source_id = $1",
		i.SourceID).Scan(&i.ID, &i.SourceID, &i.Date, &i.FileName, &i.Sha1, &line, &source, &wfstatus)

	json.Unmarshal(line, &i.Line)
	json.Unmarshal(source, &i.Source)
	json.Unmarshal(wfstatus, &i.Wfstatus)

	if err != nil {
		return err
	}

	return err
}

func (i *Source) PostSourceID(db *sql.DB) error {
	line, _ := json.Marshal(i.Line)
	source, _ := json.Marshal(i.Source)
	wfstatus, _ := json.Marshal(i.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO source(source_id, date, file_name, sha1, line, source, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (source_id) DO UPDATE SET (source_id, date, file_name, sha1, line, source, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE source.source_id = $1 RETURNING id",
		i.SourceID, i.Date, i.FileName, i.Sha1, line, source, wfstatus).Scan(&i.ID)

	if err != nil {
		return err
	}

	return nil
}

func (i *Source) PostSourceJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE source SET ` + key + ` = $2 WHERE source_id=$1;`
	//sqlStatement := `UPDATE source SET wfstatus = wtstatus || '{"`+key+`": $2}' WHERE source_id=$1;`
	_, err := db.Exec(sqlStatement, i.SourceID, v)

	return err
}

func (i *Source) PostSourceValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE source SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE source_id=$1",
		i.SourceID, value, key)

	return err
}

func (i *Source) DeleteSourceID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM source WHERE source_id=$1", i.SourceID)

	return err
}
