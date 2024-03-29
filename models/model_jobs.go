package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Jobs struct {
	ID       int                    `json:"id"`
	JobID    string                 `json:"job_id"`
	Date     string                 `json:"date"`
	FileName string                 `json:"file_name"`
	JobName  string                 `json:"job_name"`
	JobType  string                 `json:"job_type"`
	Parent   map[string]interface{} `json:"parent"`
	Line     map[string]interface{} `json:"line"`
	Original map[string]interface{} `json:"original"`
	Proxy    map[string]interface{} `json:"proxy"`
	Product  map[string]interface{} `json:"product"`
	Wfstatus map[string]interface{} `json:"wfstatus"`
}

func FindJob(db *sql.DB, values url.Values) ([]Jobs, error) {
	var where []string
	sqlStatement := `SELECT * FROM jobs WHERE wfstatus['removed'] = 'false'`
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
		if k == "doers" {
			where = append(where, fmt.Sprintf(`parent['%s'] ? '%s'`, k, v[0]))
			continue
		}
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, k, v[0]))
	}

	if len(where) > 0 {
		sqlStatement = sqlStatement + ` AND ` + strings.Join(where, " AND ")
	}

	sqlStatement = sqlStatement + fmt.Sprintf(` ORDER BY job_id DESC LIMIT %s OFFSET %s`, limit, offset)

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Jobs{}

	for rows.Next() {
		var t Jobs
		var parent, line, original, proxy, product, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(product, &t.Product)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindJobByJSON(db *sql.DB, ep string, key string, value string) ([]Jobs, error) {

	sqlStatement := fmt.Sprintf(`SELECT * FROM jobs WHERE %s['%s'] = '"%s"' ORDER BY job_id;`, ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Jobs{}

	for rows.Next() {
		var t Jobs
		var parent, line, original, proxy, product, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(product, &t.Product)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func FindJobBySHA1(db *sql.DB, value string) ([]Jobs, error) {

	sqlStatement := fmt.Sprintf("SELECT id, job_id, date, file_name, ,job_name, job_type, parent, line, original, proxy, product, wfstatus FROM jobs WHERE original->'format'->>'sha1' = '%s' ORDER BY job_id;", value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Jobs{}

	for rows.Next() {
		var t Jobs
		var parent, line, original, proxy, product, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(product, &t.Product)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetListJobs(db *sql.DB, start, count int) ([]Jobs, error) {
	rows, err := db.Query(
		"SELECT * FROM jobs ORDER BY job_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Jobs{}

	for rows.Next() {
		var t Jobs
		var parent, line, original, proxy, product, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(product, &t.Product)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func GetActiveJobs(db *sql.DB) ([]Jobs, error) {
	rows, err := db.Query(
		"SELECT * FROM jobs WHERE wfstatus ->> 'removed' = 'false' ORDER BY job_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Jobs{}

	for rows.Next() {
		var t Jobs
		var parent, line, original, proxy, product, wfstatus []byte
		if err := rows.Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(original, &t.Original)
		json.Unmarshal(proxy, &t.Proxy)
		json.Unmarshal(product, &t.Product)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Jobs) GetJobID(db *sql.DB) error {
	var parent, line, original, proxy, product, wfstatus []byte

	err := db.QueryRow("SELECT id, job_id, date, file_name, job_name, job_type, parent, line, original, proxy, product, wfstatus FROM jobs WHERE job_id = $1",
		t.JobID).Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(original, &t.Original)
	json.Unmarshal(proxy, &t.Proxy)
	json.Unmarshal(product, &t.Product)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Jobs) GetJobByID(db *sql.DB) error {
	var parent, line, original, proxy, product, wfstatus []byte

	err := db.QueryRow("SELECT id, job_id, date, file_name, job_name, job_type, parent, line, original, proxy, product, wfstatus FROM jobs WHERE id = $1",
		t.ID).Scan(&t.ID, &t.JobID, &t.Date, &t.FileName, &t.JobName, &t.JobType, &parent, &line, &original, &proxy, &product, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(original, &t.Original)
	json.Unmarshal(proxy, &t.Proxy)
	json.Unmarshal(product, &t.Product)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Jobs) PostJobID(db *sql.DB) error {
	parent, _ := json.Marshal(t.Parent)
	line, _ := json.Marshal(t.Line)
	original, _ := json.Marshal(t.Original)
	proxy, _ := json.Marshal(t.Proxy)
	product, _ := json.Marshal(t.Product)
	wfstatus, _ := json.Marshal(t.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO jobs(job_id, date, file_name, job_name, job_type, parent, line, original, proxy, product, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (job_id) DO UPDATE SET (job_id, date, file_name, job_name, job_type, parent, line, original, proxy, product, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE jobs.job_id = $1 RETURNING id",
		t.JobID, t.Date, t.FileName, t.JobName, t.JobType, parent, line, original, proxy, product, wfstatus).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *Jobs) PostJobJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE jobs SET ` + key + ` = $2 WHERE job_id=$1;`
	_, err := db.Exec(sqlStatement, t.JobID, v)

	return err
}

func (t *Jobs) PostJobValue(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE jobs SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE job_id=$1",
		t.JobID, value, key)

	return err
}

func (t *Jobs) DeleteJobID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM jobs WHERE job_id=$1", t.JobID)

	return err
}
