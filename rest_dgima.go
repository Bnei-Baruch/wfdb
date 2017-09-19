// rest_dgima.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findDgima(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findDgima(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getDgima(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getDgima(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getFilesToDgima(w http.ResponseWriter, r *http.Request) {

	files, err := getFilesToDgima(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getDgimaID(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	if err := d.getDgimaID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) getDgimaByID(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.ID, _ = strconv.Atoi(vars["id"])

	if err := d.getDgimaByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) postDgimaID(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	dc := json.NewDecoder(r.Body)
	if err := dc.Decode(&d); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := d.postDgimaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postDgimaJSON(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	dc := json.NewDecoder(r.Body)
	if err := dc.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := d.postDgimaJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postDgimaValue(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := d.postDgimaValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}


func (a *App) deleteDgimaID(w http.ResponseWriter, r *http.Request) {
	var d dgima
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	if err := d.deleteDgimaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}