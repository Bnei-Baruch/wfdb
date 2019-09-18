// rest_files.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getFile(w http.ResponseWriter, r *http.Request) {
	var k files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	if err := k.getFile(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, k)
}

func (a *App) postFile(w http.ResponseWriter, r *http.Request) {
	var k files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&k); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := k.postFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postFileValue(w http.ResponseWriter, r *http.Request) {
	var s files
	vars := mux.Vars(r)
	s.FileID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	body := r.FormValue("body")

	if value == "" {
		if err := s.postFileValue(a.DB, body, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := s.postFileStatus(a.DB, value, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postFileJSON(w http.ResponseWriter, r *http.Request) {
	var s files
	vars := mux.Vars(r)
	s.FileID = vars["id"]
	key := vars["jsonb"]
	var value map[string]interface{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	if err := s.postFileJSON(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteFile(w http.ResponseWriter, r *http.Request) {
	var k files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	if err := k.deleteFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
