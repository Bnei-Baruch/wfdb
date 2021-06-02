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

func (a *App) FindDgima(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindDgima(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) FindDgimaByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if ep == "sha1" {

		files, err := models.FindDgimaBySHA1(a.DB, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	} else {

		files, err := models.FindDgimaByJSON(a.DB, ep, key, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	}
}

func (a *App) GetDgima(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetDgima(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetFilesToDgima(w http.ResponseWriter, r *http.Request) {

	files, err := models.GetFilesToDgima(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCassetteFiles(w http.ResponseWriter, r *http.Request) {

	files, err := models.GetCassetteFiles(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetDgimaBySource(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	v := vars["id"]

	files, err := models.GetDgimaBySource(a.DB, v)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetDgimaID(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	if err := d.GetDgimaID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) GetDgimaByID(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.ID, _ = strconv.Atoi(vars["id"])

	if err := d.GetDgimaByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, d)
}

func (a *App) PostDgimaID(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	dc := json.NewDecoder(r.Body)
	if err := dc.Decode(&d); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := d.PostDgimaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.SendMessage("drim")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostDgimaJSON(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	dc := json.NewDecoder(r.Body)
	if err := dc.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := d.PostDgimaJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.SendMessage("drim")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostDgimaValue(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := d.PostDgimaValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.SendMessage("drim")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteDgimaID(w http.ResponseWriter, r *http.Request) {
	var d models.Dgims
	vars := mux.Vars(r)
	d.DgimaID = vars["id"]

	if err := d.DeleteDgimaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go a.SendMessage("drim")

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
