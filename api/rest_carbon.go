// rest_carbon.go

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

func (a *App) FindCarbonFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindCarbonFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCarbonFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetCarbonFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c models.Carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	if err := c.GetCarbonFile(a.DB); err != nil {
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

func (a *App) PostCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c models.Carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := c.PostCarbonFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteCarbonFile(w http.ResponseWriter, r *http.Request) {
	var c models.Carbon
	vars := mux.Vars(r)
	c.CarbonID = vars["id"]

	if err := c.DeleteCarbonFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
