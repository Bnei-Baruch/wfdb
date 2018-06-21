// rest_metus.go

package main

import (
	"net/http"
	_ "github.com/denisenkom/go-mssqldb"
	"fmt"
)

func (a *App) findMetus(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	fmt.Println("  value db:", value)

	files, err := findMetus(a.MSDB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, files)
}