// rest_metus.go

package api

import (
	"database/sql"
	"github.com/Bnei-Baruch/wfdb/models"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (a *App) FindMetus(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	//fmt.Println("  value db:", value)

	files, err := models.FindMetus(a.MSDB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) GetMetusByID(w http.ResponseWriter, r *http.Request) {
	var c models.Metus
	vars := mux.Vars(r)
	c.MetusID, _ = strconv.Atoi(vars["id"])

	if err := c.GetMetusMeta(a.MSDB, c.MetusID, "sha1"); err != nil {
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
