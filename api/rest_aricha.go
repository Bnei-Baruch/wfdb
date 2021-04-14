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

func (a *App) FindAricha(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindAricha(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) FindArichaByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if ep == "sha1" {

		files, err := models.FindArichaBySHA1(a.DB, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	} else {

		files, err := models.FindArichaByJSON(a.DB, ep, key, value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, files)

	}

}

func (a *App) GetAricha(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetAricha(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetBdika(w http.ResponseWriter, r *http.Request) {

	files, err := models.GetBdika(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetArichaID(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ArichaID = vars["id"]

	if err := t.GetArichaID(a.DB); err != nil {
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

func (a *App) GetArichaByID(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ID, _ = strconv.Atoi(vars["id"])

	if err := t.GetArichaByID(a.DB); err != nil {
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

func (a *App) PostArichaID(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ArichaID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostArichaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostArichaJSON(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ArichaID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostArichaJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostArichaValue(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ArichaID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := t.PostArichaValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteArichaID(w http.ResponseWriter, r *http.Request) {
	var t models.Aricha
	vars := mux.Vars(r)
	t.ArichaID = vars["id"]

	if err := t.DeleteArichaID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
