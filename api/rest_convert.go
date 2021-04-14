// rest_convert.go

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

func (a *App) FindConvert(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindConvert(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) FindConvertByJSON(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindConvertByJSON(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetConvert(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetConvert(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetConvertByID(w http.ResponseWriter, r *http.Request) {
	var i models.Convert
	vars := mux.Vars(r)
	i.ConvertID = vars["id"]

	if err := i.GetConvertByID(a.DB); err != nil {
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

func (a *App) PostConvert(w http.ResponseWriter, r *http.Request) {
	var i models.Convert
	vars := mux.Vars(r)
	i.ConvertID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := i.PostConvert(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostConvertValue(w http.ResponseWriter, r *http.Request) {
	var i models.Convert
	vars := mux.Vars(r)
	i.ConvertID = vars["id"]
	key := vars["key"]
	value := r.FormValue("value")

	if err := i.PostConvertValue(a.DB, key, value); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostConvertJSON(w http.ResponseWriter, r *http.Request) {
	var i models.Convert
	vars := mux.Vars(r)
	i.ConvertID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := i.PostConvertJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteConvert(w http.ResponseWriter, r *http.Request) {
	var i models.Convert
	vars := mux.Vars(r)
	i.ConvertID = vars["id"]

	if err := i.DeleteConvert(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
