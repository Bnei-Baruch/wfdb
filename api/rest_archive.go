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

func (a *App) TestAr() {

}

func (a *App) FindArFiles(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	arfiles, err := models.FindArFiles(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, arfiles)
}

func (a *App) GetArFiles(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 {
		count = 100
	}
	if start < 0 {
		start = 0
	}

	arfiles, err := models.GetArFiles(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, arfiles)
}

func (a *App) GetArFile(w http.ResponseWriter, r *http.Request) {
	var arch models.Archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	if err := arch.GetArFile(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, arch)
}

func (a *App) PostArFile(w http.ResponseWriter, r *http.Request) {
	var arch models.Archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&arch); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := arch.PostArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) UpdateArFile(w http.ResponseWriter, r *http.Request) {
	var arch models.Archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&arch); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := arch.UpdateArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteArFile(w http.ResponseWriter, r *http.Request) {
	var arch models.Archive
	vars := mux.Vars(r)
	arch.ArchiveID = vars["id"]

	if err := arch.DeleteArFile(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
