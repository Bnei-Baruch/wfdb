// rest_capture.go

package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) findCaptures(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	captures, err := findCaptures(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, captures)
}

func (a *App) getCaptures(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	captures, err := getCaptures(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, captures)
}

func (a *App) getCapture(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.getCapture(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c.Data)
}

func (a *App) postCapture(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.postCapture(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) updateCapture(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.updateCapture(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteCapture(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.deleteCapture(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}