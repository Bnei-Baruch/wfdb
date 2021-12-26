package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Users struct {
	ID        int         `json:"id"`
	UserID    string      `json:"user_id"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     string      `json:"email"`
	Props     interface{} `json:"properties"`
}

func FindUsers(db *sql.DB, values url.Values) ([]Users, error) {

	var where []string
	sqlStatement := `SELECT * FROM users`
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
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, k, v[0]))
	}

	if len(where) > 0 {
		sqlStatement = sqlStatement + ` AND ` + strings.Join(where, " AND ")
	}

	sqlStatement = sqlStatement + fmt.Sprintf(` ORDER BY oid DESC LIMIT %s OFFSET %s`, limit, offset)

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Users{}

	for rows.Next() {
		var t Users
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindUsersByJSON(db *sql.DB, ep string, key string, value string) ([]Users, error) {

	sqlStatement := fmt.Sprintf(`SELECT * FROM users WHERE %s['%s'] = '"%s"' ORDER BY oid;`, ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Users{}

	for rows.Next() {
		var t Users
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetListUsers(db *sql.DB, start, count int) ([]Users, error) {
	rows, err := db.Query(
		"SELECT * FROM users ORDER BY user_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Users{}

	for rows.Next() {
		var t Users
		var properties []byte
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties); err != nil {
			return nil, err
		}
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Users) GetUsersID(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT * FROM users WHERE user_id = $1",
		t.UserID).Scan(
		&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Users) GetUsersByID(db *sql.DB) error {
	var properties []byte

	err := db.QueryRow("SELECT * FROM users WHERE id = $1",
		t.ID).Scan(&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Users) PostUsersID(db *sql.DB) error {
	properties, _ := json.Marshal(t.Props)

	err := db.QueryRow(
		"INSERT INTO users(user_id, first_name, last_name, email, properties) VALUES($1, $2, $3, $4, $5) ON CONFLICT (user_id) DO UPDATE SET (user_id, first_name, last_name, email, properties) = ($1, $2, $3, $4, $5) WHERE users.user_id = $1 RETURNING id",
		&t.ID, &t.UserID, &t.FirstName, &t.LastName, &t.Email, &properties).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *Users) PostUsersJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE users SET ` + key + ` = $2 WHERE oid=$1;`
	_, err := db.Exec(sqlStatement, t.UserID, v)

	return err
}

func (t *Users) SetUsersJSON(db *sql.DB, value interface{}, key string, prop string) error {

	v, _ := json.Marshal(value)
	sqlCmd := "UPDATE users SET " + prop + " = " + prop + " || json_build_object($3::text, $2::jsonb)::jsonb WHERE user_id=$1"
	_, err := db.Exec(sqlCmd, t.UserID, v, key)

	return err
}

func (t *Users) PostUsersStatus(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE users SET properties = properties || json_build_object($3::text, $2::bool)::jsonb WHERE user_id=$1",
		t.UserID, value, key)

	return err
}

func (t *Users) PostUsersProp(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE users SET properties = properties || json_build_object($3::text, $2::text)::jsonb WHERE user_id=$1",
		t.UserID, value, key)

	return err
}

func (t *Users) DeleteUsersID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE user_id=$1", t.UserID)

	return err
}
