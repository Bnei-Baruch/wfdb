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

func (a *App) FindFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	files, err := models.FindFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	files, err := models.GetFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetActiveFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	language := vars["language"]
	product_id := r.FormValue("product_id")

	if language == "find" {
		return
	}

	files, err := models.GetActiveFiles(a.DB, language, product_id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetFile(w http.ResponseWriter, r *http.Request) {
	var k models.Files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	if err := k.GetFile(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, k)
}

func (a *App) PostFile(w http.ResponseWriter, r *http.Request) {
	var k models.Files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&k); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := k.PostFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostFileStatus(w http.ResponseWriter, r *http.Request) {
	var s models.Files
	vars := mux.Vars(r)
	s.FileID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := s.PostFileStatus(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostFileValue(w http.ResponseWriter, r *http.Request) {
	var s models.Files
	vars := mux.Vars(r)
	s.FileID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")

	if err := s.PostFileValue(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostFileJSON(w http.ResponseWriter, r *http.Request) {
	var s models.Files
	vars := mux.Vars(r)
	s.FileID = vars["id"]
	key := vars["jsonb"]
	var value map[string]interface{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	if err := s.PostFileJSON(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteFile(w http.ResponseWriter, r *http.Request) {
	var k models.Files
	vars := mux.Vars(r)
	k.FileID = vars["id"]

	if err := k.DeleteFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
