// rest_archive.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findArFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	arfiles, err := findArFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, arfiles)
}

func (a *App) getArFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	arfiles, err := getArFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, arfiles)
}

func (a *App) getArFile(w http.ResponseWriter, r *http.Request) {
	var arch archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	if err := arch.getArFile(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, arch)
}

func (a *App) postArFile(w http.ResponseWriter, r *http.Request) {
	var arch archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&arch); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := arch.postArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateArFile(w http.ResponseWriter, r *http.Request) {
	var arch archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&arch); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := arch.updateArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteArFile(w http.ResponseWriter, r *http.Request) {
	var arch archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	if err := arch.deleteArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}