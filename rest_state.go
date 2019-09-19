// rest_state.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findState(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	states, err := findStates(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, states)
}

func (a *App) getStates(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	states, err := getStates(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, states)
}

func (a *App) getState(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	if err := s.getState(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s.Data)
}

func (a *App) postState(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&s.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := s.postState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateState(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&s.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := s.updateState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postStateValue(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	body := r.FormValue("body")

	if value == "" {
		if err := s.postStateValue(a.DB, body, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := s.postStateStatus(a.DB, value, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postStateJSON(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	key := vars["jsonb"]
	var value map[string]interface{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	if err := s.postStateJSON(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteState(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	if err := s.deleteState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteStateJSON(w http.ResponseWriter, r *http.Request) {
	var s state
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	value := vars["jsonb"]

	if err := s.deleteStateJSON(a.DB, value); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
