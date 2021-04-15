package api

import (
	//"database/sql"
	//"encoding/json"
	"net/http"

	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
)

func (a *App) PostTask(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")
	file, _, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	log := string(data[:])

	fmt.Println(w, "%v", log)
	fmt.Fprintf(w, "%v", key)
	fmt.Fprintf(w, "%v", value)

	//states, err := FindStates(a.DB, key, value)
	//if err != nil {
	//	respondWithError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}

	//respondWithJSON(w, http.StatusOK, states)
	respondWithJSON(w, http.StatusCreated, "File uploaded successfully!.")
}
