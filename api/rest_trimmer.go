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

func (a *App) FindTrimmer(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindTrimmer(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) FindTrimmerByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if ep == "sha1" {

		files, err := models.FindTrimmerBySHA1(a.DB, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	} else {

		files, err := models.FindTrimmerByJSON(a.DB, ep, key, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	}
}

func (a *App) GetTrimmer(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetTrimmer(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetFilesToTrim(w http.ResponseWriter, r *http.Request) {

	files, err := models.GetFilesToTrim(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.GetTrimmerID(a.DB); err != nil {
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

func (a *App) GetTrimmerByID(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.ID, _ = strconv.Atoi(vars["id"])

	if err := t.GetTrimmerByID(a.DB); err != nil {
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

func (a *App) PostTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostTrimmerID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("trimmer")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostTrimmerJSON(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostTrimmerJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("trimmer")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostTrimmerValue(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := t.PostTrimmerValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("trimmer")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteTrimmerID(w http.ResponseWriter, r *http.Request) {
	var t models.Trimmer
	vars := mux.Vars(r)
	t.TrimID = vars["id"]

	if err := t.DeleteTrimmerID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.ReportMonitor("trimmer")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
