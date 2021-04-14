package api

import (
	"database/sql"
	"encoding/json"
	"github.com/Bnei-Baruch/wfdb/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func (a *App) FindCapture(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindCapture(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCapture(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetCapture(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCaptureID(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.GetCaptureID(a.DB); err != nil {
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

func (a *App) GetCassetteID(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
	vars := mux.Vars(r)
	c.StopName = vars["id"]

	if err := c.GetCassetteID(a.DB); err != nil {
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

func (a *App) PostCaptureID(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.PostCaptureID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostCaptureJSON(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
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

	if err := c.PostCaptureJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostCaptureValue(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	//value := strconv.FormatBool(r.FormValue("value"))
	//value, _ := strconv.ParseBool(r.FormValue("value"))

	if err := c.PostCaptureValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteCaptureID(w http.ResponseWriter, r *http.Request) {
	var c models.Capture
	vars := mux.Vars(r)
	c.CaptureID = vars["id"]

	if err := c.DeleteCaptureID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
