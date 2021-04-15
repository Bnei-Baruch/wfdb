package api

import (
	"database/sql"
	"encoding/json"
	"github.com/Bnei-Baruch/wfdb/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func (a *App) FindState(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")

	states, err := models.FindStates(a.DB, key, value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, states)
}

func (a *App) GetStateByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]

	states, err := models.GetStateByTag(a.DB, tag)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, states)
}

func (a *App) GetStates(w http.ResponseWriter, r *http.Request) {

	states, err := models.GetStates(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, states)
}

func (a *App) GetState(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	if err := s.GetState(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s.Data)
}

func (a *App) GetStateJSON(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.Tag = vars["tag"]
	s.StateID = vars["id"]
	key := vars["jsonb"]

	if err := s.GetStateJSON(a.DB, key); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s.Data)
}

func (a *App) PostState(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.Tag = vars["tag"]
	s.StateID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&s.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := s.PostState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) UpdateState(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&s.Data); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()

	if err := s.UpdateState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostStateValue(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	key := vars["jsonb"]
	value := r.FormValue("value")
	body := r.FormValue("body")

	if value == "" {
		if err := s.PostStateValue(a.DB, body, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := s.PostStateStatus(a.DB, value, key); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) PostStateJSON(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	key := vars["jsonb"]
	var value map[string]interface{}
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&value); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	if err := s.PostStateJSON(a.DB, value, key); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteState(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]

	if err := s.DeleteState(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteStateJSON(w http.ResponseWriter, r *http.Request) {
	var s models.State
	vars := mux.Vars(r)
	s.StateID = vars["id"]
	value := vars["jsonb"]

	if err := s.DeleteStateJSON(a.DB, value); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
