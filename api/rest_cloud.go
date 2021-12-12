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

func (a *App) FindCloud(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	files, err := models.FindCloud(a.DB, values)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) FindCloudByJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ep := vars["jsonb"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindCloudByJSON(a.DB, ep, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, files)

}

func (a *App) GetListClouds(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetListClouds(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetCloudID(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]

	if err := t.GetCloudID(a.DB); err != nil {
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

func (a *App) GetCloudByID(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.ID, _ = strconv.Atoi(vars["id"])

	if err := t.GetCloudByID(a.DB); err != nil {
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

func (a *App) PostCloudID(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostCloudID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostCloudJSON(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]
	key := vars["jsonb"]
	var jsonb map[string]interface{}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&jsonb); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := t.PostCloudJSON(a.DB, jsonb, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) SetCloudJSON(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]
	key := vars["jsonb"]
	prop := vars["prop"]
	var value interface{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	if err := t.SetCloudJSON(a.DB, value, key, prop); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostCloudStatus(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if err := t.PostCloudStatus(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostCloudProp(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]
	key := r.FormValue("key")
	value := r.FormValue("value")

	if err := t.PostCloudProp(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteCloudID(w http.ResponseWriter, r *http.Request) {
	var t models.Clouds
	vars := mux.Vars(r)
	t.OID = vars["id"]

	if err := t.DeleteCloudID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
