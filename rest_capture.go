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

func (a *App) findCapture(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := findCapture(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getCapture(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := getCapture(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getCaptureID(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.getCaptureID(a.DB); err != nil {
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

func (a *App) getCassetteID(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.StopName = vars["id"]

	if err := c.getCassetteID(a.DB); err != nil {
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

func (a *App) postCaptureID(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.postCaptureID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postCaptureJSON(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.postCaptureJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postCaptureValue(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	//value := strconv.FormatBool(r.FormValue("value"))
	//value, _ := strconv.ParseBool(r.FormValue("value"))

	if err := c.postCaptureValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteCaptureID(w http.ResponseWriter, r *http.Request) {
	var c capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.deleteCaptureID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}