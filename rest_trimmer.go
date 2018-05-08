// rest_trimmer.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findTrimmer(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findTrimmer(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) findTrimmerByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findTrimmerByJSON(a.DB, ep, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getTrimmer(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getTrimmer(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getFilesToTrim(w http.ResponseWriter, r *http.Request) {

	files, err := getFilesToTrim(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.getTrimmerID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) getTrimmerByID(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.ID, _ = strconv.Atoi(vars["id"])

	if err := t.getTrimmerByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) postTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.postTrimmerID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postTrimmerJSON(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.postTrimmerJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postTrimmerValue(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := t.postTrimmerValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (a *App) deleteTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.deleteTrimmerID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}