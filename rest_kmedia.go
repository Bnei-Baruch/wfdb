// rest_kmedia.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findKmFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findKmFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getKmFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getKmFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getKmFile(w http.ResponseWriter, r *http.Request) {
	var k kmedia
	vars := mux.Vars(r)
	k.KmediaID = vars["id"]

	if err := k.getKmFile(a.DB); err != nil {
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

func (a *App) postKmFile(w http.ResponseWriter, r *http.Request) {
	var k kmedia
	vars := mux.Vars(r)
	k.KmediaID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&k); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := k.postKmFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteKmFile(w http.ResponseWriter, r *http.Request) {
	var k kmedia
	vars := mux.Vars(r)
	k.KmediaID = vars["id"]

	if err := k.deleteKmFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}