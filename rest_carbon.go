// rest_carbon.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findCarbonFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findCarbonFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getCarbonFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getCarbonFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	if err := c.getCarbonFile(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) postCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.postCarbonFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	if err := c.deleteCarbonFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}