// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	MSDB   *sql.DB
}

func (a *App) Initialize(user string, password string, dbname string, host string, user_id string, pass string, name string) {
	connectionString :=
		fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	conString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable;", host, user_id, pass, name)
	a.MSDB, err = sql.Open("mssql", conString)
	if err != nil {
		fmt.Println("  Error open db:", err.Error())
	}

	//defer a.MSDB.Close()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Accept-Encoding"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/metus/find", a.findMetus).Methods("GET")
	a.Router.HandleFunc("/metus/{id:[0-9]+}", a.getMetusByID).Methods("GET")
	// Capture
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.postCaptureID).Methods("PUT")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}/wfstatus/{jsonb}", a.postCaptureValue).Methods("POST")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}/{jsonb}", a.postCaptureJSON).Methods("POST")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.getCaptureID).Methods("GET")
	a.Router.HandleFunc("/capture/{id:[0-9]+}", a.getCassetteID).Methods("GET")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.deleteCaptureID).Methods("DELETE")
	a.Router.HandleFunc("/capture", a.getCapture).Methods("GET")
	a.Router.HandleFunc("/capture/find", a.findCapture).Methods("GET")
	// Archive
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.postArFile).Methods("PUT")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.updateArFile).Methods("POST")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.getArFile).Methods("GET")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.deleteArFile).Methods("DELETE")
	a.Router.HandleFunc("/archive", a.getArFiles).Methods("GET")
	a.Router.HandleFunc("/archive/find", a.findArFiles).Methods("GET")
	// Convert
	a.Router.HandleFunc("/convert", a.getConvert).Methods("GET")
	a.Router.HandleFunc("/convert/find", a.findConvert).Methods("GET")
	a.Router.HandleFunc("/convert/langcheck", a.findConvertByJSON).Methods("GET")
	a.Router.HandleFunc("/convert/{id:[t|d|a][0-9]+}", a.getConvertByID).Methods("GET")
	a.Router.HandleFunc("/convert/{id}", a.postConvert).Methods("PUT")
	a.Router.HandleFunc("/convert/{id}/{key}", a.postConvertValue).Methods("POST")
	a.Router.HandleFunc("/convert/{id}/{jsonb}", a.postConvertJSON).Methods("PUT")
	a.Router.HandleFunc("/convert/{id}", a.deleteConvert).Methods("DELETE")
	// Carbon
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.postCarbonFile).Methods("PUT")
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.getCarbonFile).Methods("GET")
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.deleteCarbonFile).Methods("DELETE")
	a.Router.HandleFunc("/carbon", a.getCarbonFiles).Methods("GET")
	a.Router.HandleFunc("/carbon/find", a.findCarbonFiles).Methods("GET")
	// Kmedia
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.postKmFile).Methods("PUT")
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.getKmFile).Methods("GET")
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.deleteKmFile).Methods("DELETE")
	a.Router.HandleFunc("/kmedia", a.getKmFiles).Methods("GET")
	a.Router.HandleFunc("/kmedia/find", a.findKmFiles).Methods("GET")
	// Insert
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.postInsertFile).Methods("PUT")
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.getInsertFile).Methods("GET")
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.deleteInsertFile).Methods("DELETE")
	a.Router.HandleFunc("/insert", a.getInsertFiles).Methods("GET")
	a.Router.HandleFunc("/insert/find", a.findInsertFiles).Methods("GET")
	a.Router.HandleFunc("/insert/line", a.findInsertByJSON).Methods("GET")
	// Ingest
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.postIngestID).Methods("PUT")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}/wfstatus/{jsonb}", a.postIngestValue).Methods("POST")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}/{jsonb}", a.postIngestJSON).Methods("POST")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.getIngestID).Methods("GET")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.deleteIngestID).Methods("DELETE")
	a.Router.HandleFunc("/ingest", a.getIngest).Methods("GET")
	a.Router.HandleFunc("/ingest/find", a.findIngest).Methods("GET")
	// Dgima
	a.Router.HandleFunc("/drim", a.getFilesToDgima).Methods("GET")
	a.Router.HandleFunc("/drim/{id}", a.getDgimaBySource).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.postDgimaID).Methods("PUT")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}/wfstatus/{jsonb}", a.postDgimaValue).Methods("POST")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}/{jsonb}", a.postDgimaJSON).Methods("POST")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.getDgimaID).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:[0-9]+}", a.getDgimaByID).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.deleteDgimaID).Methods("DELETE")
	a.Router.HandleFunc("/dgima", a.getDgima).Methods("GET")
	a.Router.HandleFunc("/dgima/find", a.findDgima).Methods("GET")
	a.Router.HandleFunc("/dgima/{jsonb}", a.findDgimaByJSON).Methods("GET")
	// Trimmer
	a.Router.HandleFunc("/trim", a.getFilesToTrim).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.postTrimmerID).Methods("PUT")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}/wfstatus/{jsonb}", a.postTrimmerValue).Methods("POST")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}/{jsonb}", a.postTrimmerJSON).Methods("POST")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.getTrimmerID).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:[0-9]+}", a.getTrimmerByID).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.deleteTrimmerID).Methods("DELETE")
	a.Router.HandleFunc("/trimmer", a.getTrimmer).Methods("GET")
	a.Router.HandleFunc("/trimmer/find", a.findTrimmer).Methods("GET")
	a.Router.HandleFunc("/trimmer/{jsonb}", a.findTrimmerByJSON).Methods("GET")
	// Aricha
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.postArichaID).Methods("PUT")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}/wfstatus/{jsonb}", a.postArichaValue).Methods("POST")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}/{jsonb}", a.postArichaJSON).Methods("POST")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.getArichaID).Methods("GET")
	a.Router.HandleFunc("/aricha/{id:[0-9]+}", a.getArichaByID).Methods("GET")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.deleteArichaID).Methods("DELETE")
	a.Router.HandleFunc("/aricha", a.getAricha).Methods("GET")
	a.Router.HandleFunc("/bdika", a.getBdika).Methods("GET")
	a.Router.HandleFunc("/aricha/find", a.findAricha).Methods("GET")
	a.Router.HandleFunc("/aricha/{jsonb}", a.findArichaByJSON).Methods("GET")
	// Jobs
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.postJobID).Methods("PUT")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}/wfstatus/{jsonb}", a.postJobValue).Methods("POST")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}/{jsonb}", a.postJobJSON).Methods("POST")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.getJobID).Methods("GET")
	a.Router.HandleFunc("/jobs/{id:[0-9]+}", a.getJobByID).Methods("GET")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.deleteJobID).Methods("DELETE")
	a.Router.HandleFunc("/jobs_list", a.getListJobs).Methods("GET")
	a.Router.HandleFunc("/jobs", a.getActiveJobs).Methods("GET")
	a.Router.HandleFunc("/jobs/find", a.findJob).Methods("GET")
	a.Router.HandleFunc("/jobs/{jsonb}", a.findJobByJSON).Methods("GET")
	// Files
	a.Router.HandleFunc("/files/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.postFile).Methods("PUT")
	a.Router.HandleFunc("/files/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.getFile).Methods("GET")
	a.Router.HandleFunc("/files/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.deleteFile).Methods("DELETE")
	a.Router.HandleFunc("/files/{id:[a-z0-9_-]+\\.[a-z0-9]+}/line/{jsonb}", a.postFileJSON).Methods("PUT")
	a.Router.HandleFunc("/files/{id:[a-z0-9_-]+\\.[a-z0-9]+}/line/{jsonb}", a.postFileValue).Methods("POST")
	a.Router.HandleFunc("/files", a.getFiles).Methods("GET")
	a.Router.HandleFunc("/files/find", a.findFiles).Methods("GET")
	// Tasks
	a.Router.HandleFunc("/task", a.postTask).Methods("POST")
	// Labels
	a.Router.HandleFunc("/labels", a.getLabels).Methods("GET")
	a.Router.HandleFunc("/label/{id:[0-9]+}", a.getLabel).Methods("GET")
	a.Router.HandleFunc("/labels/find", a.findLabels).Methods("GET")
	// State
	a.Router.HandleFunc("/states", a.getStates).Methods("GET")
	a.Router.HandleFunc("/{tag}", a.getStateByTag).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}", a.getState).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.getStateJSON).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}", a.postState).Methods("PUT")
	a.Router.HandleFunc("/{tag}/{id}", a.updateState).Methods("POST")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.postStateJSON).Methods("PUT")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.postStateValue).Methods("POST")
	a.Router.HandleFunc("/{tag}/{id}", a.deleteState).Methods("DELETE")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.deleteStateJSON).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}
