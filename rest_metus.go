// rest_metus.go

package main

import (
	"net/http"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"database/sql"
	"strconv"
)

func (a *App) findMetus(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	//fmt.Println("  value db:", value)

	files, err := findMetus(a.MSDB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}

func (a *App) getMetusByID(w http.ResponseWriter, r *http.Request) {
	var c metus
	vars := mux.Vars(r)
	c.MetusID, _ = strconv.Atoi(vars["id"])

	if err := c.getMetusMeta(a.MSDB, c.MetusID, "sha1"); err != nil {
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