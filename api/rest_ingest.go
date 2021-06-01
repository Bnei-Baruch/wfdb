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

func (a *App) FindIngest(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindIngest(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetIngest(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetIngest(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetIngestID(w http.ResponseWriter, r *http.Request) {
	var i models.Ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	if err := i.GetIngestID(a.DB); err != nil {
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

func (a *App) PostIngestID(w http.ResponseWriter, r *http.Request) {
	var i models.Ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := i.PostIngestID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("ingest")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostIngestJSON(w http.ResponseWriter, r *http.Request) {
	var i models.Ingest
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

	if err := i.PostIngestJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("ingest")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostIngestValue(w http.ResponseWriter, r *http.Request) {
	var i models.Ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	//value := strconv.FormatBool(r.FormValue("value"))
	//value, _ := strconv.ParseBool(r.FormValue("value"))

	if err := i.PostIngestValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("ingest")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteIngestID(w http.ResponseWriter, r *http.Request) {
	var i models.Ingest
	vars := mux.Vars(r)
	i.CaptureID = vars["id"]

	if err := i.DeleteIngestID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("ingest")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
