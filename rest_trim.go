// rest_trim.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findTrimes(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	trimes, err := findTrimes(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, trimes)
}

func (a *App) getTrimes(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	trimes, err := getTrimes(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, trimes)
}

func (a *App) getTrim(w http.ResponseWriter, r *http.Request) {
	var t trim
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.getTrim(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t.Data)
}

func (a *App) postTrim(w http.ResponseWriter, r *http.Request) {
	var t trim
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.postTrim(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateTrim(w http.ResponseWriter, r *http.Request) {
	var t trim
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.updateTrim(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteTrim(w http.ResponseWriter, r *http.Request) {
	var t trim
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.deleteTrim(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}