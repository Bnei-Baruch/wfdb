// model_state.go

package main

import (
	"database/sql"
	"encoding/json"
)

type state struct {
	ID      int                    `json:"id"`
	StateID string                 `json:"state_id"`
	Data    map[string]interface{} `json:"data"`
}

func findStates(db *sql.DB, key string, value string) ([]state, error) {
	rows, err := db.Query(
		"SELECT id, state_id, data FROM state WHERE data @> json_build_object($1::text, $2::text)::jsonb",
		key, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	states := []state{}

	for rows.Next() {
		var s state
		var obj []byte
		if err := rows.Scan(&s.ID, &s.StateID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &s.Data)
		states = append(states, s)
	}

	return states, nil
}

func getStates(db *sql.DB, start, count int) ([]state, error) {
	rows, err := db.Query(
		"SELECT id, state_id, data FROM state ORDER BY state_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	states := []state{}

	for rows.Next() {
		var s state
		var obj []byte
		if err := rows.Scan(&s.ID, &s.StateID, &obj); err != nil {
			return nil, err
		}
		json.Unmarshal(obj, &s.Data)
		states = append(states, s)
	}

	return states, nil
}

func (s *state) getState(db *sql.DB) error {
	var obj []byte
	err := db.QueryRow("SELECT data FROM state WHERE state_id = $1",
		s.StateID).Scan(&obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, &s.Data)

	return err
}

func (s *state) postState(db *sql.DB) error {
	v, _ := json.Marshal(s.Data)

	err := db.QueryRow(
		"INSERT INTO state(state_id, data) VALUES($1, $2) ON CONFLICT (state_id) DO UPDATE SET data = $2 WHERE state.state_id = $1 RETURNING id",
		s.StateID, v).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *state) updateState(db *sql.DB) error {
	v, _ := json.Marshal(s.Data)
	_, err :=
		db.Exec("UPDATE state SET data=$2 WHERE state_id=$1",
			s.StateID, v)

	return err
}

func (s *state) postStateStatus(db *sql.DB, value, key string) error {
	_, err := db.Exec("UPDATE state SET data = data || json_build_object($3::text, $2::bool)::jsonb WHERE state_id=$1",
		s.StateID, value, key)

	return err
}

func (s *state) postStateValue(db *sql.DB, value string, key string) error {
	_, err := db.Exec("UPDATE state SET data = data || json_build_object($3::text, $2::text)::jsonb WHERE state_id=$1",
		s.StateID, value, key)

	return err
}

func (s *state) postStateJSON(db *sql.DB, value interface{}, key string) error {
	v, _ := json.Marshal(value)
	_, err := db.Exec("UPDATE state SET data = data || json_build_object($3::text, $2::jsonb)::jsonb WHERE state_id=$1",
		s.StateID, v, key)

	return err
}

func (s *state) deleteState(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM state WHERE state_id=$1", s.StateID)

	return err
}
