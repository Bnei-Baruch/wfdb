// rest_jobs.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findJob(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findJob(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) findJobByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if(ep == "sha1") {

		files, err := findJobBySHA1(a.DB, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	} else {

		files, err := findJobByJSON(a.DB, ep, key, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	}

}

func (a *App) getListJobs(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getListJobs(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getActiveJobs(w http.ResponseWriter, r *http.Request) {

	files, err := getActiveJobs(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getJobID(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.JobID = vars["id"]

	if err := t.getJobID(a.DB); err != nil {
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

func (a *App) getJobByID(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.ID, _ = strconv.Atoi(vars["id"])

	if err := t.getJobByID(a.DB); err != nil {
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

func (a *App) postJobID(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.JobID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.postJobID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postJobJSON(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.JobID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.postJobJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postJobValue(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.JobID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := t.postJobValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteJobID(w http.ResponseWriter, r *http.Request) {
	var t jobs
	vars := mux.Vars(r)
	t.JobID = vars["id"]

	if err := t.deleteJobID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}