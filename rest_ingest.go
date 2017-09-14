// rest_ingest.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findIngest(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findIngest(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getIngest(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getIngest(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getIngestID(w http.ResponseWriter, r *http.Request) {
	var i ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	if err := i.getIngestID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

func (a *App) postIngestID(w http.ResponseWriter, r *http.Request) {
	var i ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := i.postIngestID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postIngestJSON(w http.ResponseWriter, r *http.Request) {
	var i ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := i.postIngestJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postIngestValue(w http.ResponseWriter, r *http.Request) {
	var i ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	//value := strconv.FormatBool(r.FormValue("value"))
	//value, _ := strconv.ParseBool(r.FormValue("value"))

	if err := i.postIngestValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteIngestID(w http.ResponseWriter, r *http.Request) {
	var i ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	if err := i.deleteIngestID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}