package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/wfdb/common"
	"github.com/Bnei-Baruch/wfdb/pkg/middleware"
	"github.com/coreos/go-oidc"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type App struct {
	Router        *mux.Router
	Handler       http.Handler
	DB            *sql.DB
	MSDB          *sql.DB
	tokenVerifier *oidc.IDTokenVerifier
	Msg           mqtt.Client
}

func (a *App) InitDB() {
	user := common.APP_DB_USERNAME
	password := common.APP_DB_PASSWORD
	dbname := common.APP_DB_NAME
	dbhost := common.APP_DB_HOST
	dbport := common.APP_DB_PORT
	host := common.METUS_DB_HOST
	user_id := common.METUS_DB_USERNAME
	pass := common.METUS_DB_PASSWORD
	name := common.METUS_DB_NAME

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, dbhost, dbport, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("Error open wfdb")
	}

	conString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable;", host, user_id, pass, name)
	a.MSDB, err = sql.Open("mssql", conString)
	if err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("Error open metus db")
	}
}

func (a *App) Initialize(accountsUrl string, skipAuth bool) {
	middleware.InitLog()
	log.Info().Str("source", "APP").Msg("initializing app")
	a.InitApp(accountsUrl, skipAuth)
}

func (a *App) InitApp(accountsUrl string, skipAuth bool) {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.initMQTT()

	if !skipAuth {
		a.initOidc(accountsUrl)
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
	})

	a.Handler = middleware.ContextMiddleware(
		middleware.LoggingMiddleware(
			middleware.RecoveryMiddleware(
				middleware.RealIPMiddleware(
					corsMiddleware.Handler(
						middleware.AuthenticationMiddleware(a.tokenVerifier, skipAuth)(
							a.Router))))))
}

func (a *App) initOidc(issuer string) {
	oidcProvider, err := oidc.NewProvider(context.TODO(), issuer)
	if err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("oidc.NewProvider")
	}

	a.tokenVerifier = oidcProvider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
	})
}

func (a *App) Run(listenAddr string) {
	addr := listenAddr
	if addr == "" {
		addr = ":8080"
	}

	log.Info().Str("source", "APP").Msgf("app run %s", addr)
	if err := http.ListenAndServe(addr, a.Handler); err != nil {
		log.Fatal().Str("source", "APP").Err(err).Msg("http.ListenAndServe")
	}
}

func (a *App) initializeRoutes() {
	//a.Router.Use(a.loggingMiddleware)
	a.Router.HandleFunc("/metus/find", a.FindMetus).Methods("GET")
	a.Router.HandleFunc("/metus/{id:[0-9]+}", a.GetMetusByID).Methods("GET")
	// Capture
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.PostCaptureID).Methods("PUT")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}/wfstatus/{jsonb}", a.PostCaptureValue).Methods("POST")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}/{jsonb}", a.PostCaptureJSON).Methods("POST")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.GetCaptureID).Methods("GET")
	a.Router.HandleFunc("/capture/{id:[0-9]+}", a.GetCassetteID).Methods("GET")
	a.Router.HandleFunc("/capture/{id:c[0-9]+}", a.DeleteCaptureID).Methods("DELETE")
	a.Router.HandleFunc("/capture", a.GetCapture).Methods("GET")
	a.Router.HandleFunc("/capture/find", a.FindCapture).Methods("GET")
	// Archive
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.PostArFile).Methods("PUT")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.UpdateArFile).Methods("POST")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.GetArFile).Methods("GET")
	a.Router.HandleFunc("/archive/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.DeleteArFile).Methods("DELETE")
	a.Router.HandleFunc("/archive", a.GetArFiles).Methods("GET")
	a.Router.HandleFunc("/archive/find", a.FindArFiles).Methods("GET")
	// Convert
	a.Router.HandleFunc("/convert", a.GetConvert).Methods("GET")
	a.Router.HandleFunc("/convert/find", a.FindConvert).Methods("GET")
	a.Router.HandleFunc("/convert/langcheck", a.FindConvertByJSON).Methods("GET")
	a.Router.HandleFunc("/convert/{id:[t|d|a][0-9]+}", a.GetConvertByID).Methods("GET")
	a.Router.HandleFunc("/convert/{id}", a.PostConvert).Methods("PUT")
	a.Router.HandleFunc("/convert/{id}/{key}", a.PostConvertValue).Methods("POST")
	a.Router.HandleFunc("/convert/{id}/{jsonb}", a.PostConvertJSON).Methods("PUT")
	a.Router.HandleFunc("/convert/{id}", a.DeleteConvert).Methods("DELETE")
	// Carbon
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.PostCarbonFile).Methods("PUT")
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.GetCarbonFile).Methods("GET")
	a.Router.HandleFunc("/carbon/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.DeleteCarbonFile).Methods("DELETE")
	a.Router.HandleFunc("/carbon", a.GetCarbonFiles).Methods("GET")
	a.Router.HandleFunc("/carbon/find", a.FindCarbonFiles).Methods("GET")
	// Kmedia
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.PostKmFile).Methods("PUT")
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.GetKmFile).Methods("GET")
	a.Router.HandleFunc("/kmedia/{id:[a-z0-9_-]+\\.[a-z0-9]+}", a.DeleteKmFile).Methods("DELETE")
	a.Router.HandleFunc("/kmedia", a.GetKmFiles).Methods("GET")
	a.Router.HandleFunc("/kmedia/find", a.FindKmFiles).Methods("GET")
	// Insert
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.PostInsertFile).Methods("PUT")
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.GetInsertFile).Methods("GET")
	a.Router.HandleFunc("/insert/{id:i[0-9]+}", a.DeleteInsertFile).Methods("DELETE")
	a.Router.HandleFunc("/insert", a.GetInsertFiles).Methods("GET")
	a.Router.HandleFunc("/insert/find", a.FindInsertFiles).Methods("GET")
	a.Router.HandleFunc("/insert/line", a.FindInsertByJSON).Methods("GET")
	// Ingest
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.PostIngestID).Methods("PUT")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}/wfstatus/{jsonb}", a.PostIngestValue).Methods("POST")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}/{jsonb}", a.PostIngestJSON).Methods("POST")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.GetIngestID).Methods("GET")
	a.Router.HandleFunc("/ingest/{id:c[0-9]+}", a.DeleteIngestID).Methods("DELETE")
	a.Router.HandleFunc("/ingest", a.GetIngest).Methods("GET")
	a.Router.HandleFunc("/ingest/find", a.FindIngest).Methods("GET")
	// Source
	a.Router.HandleFunc("/source/{id:s[0-9]+}", a.PostSourceID).Methods("PUT")
	a.Router.HandleFunc("/source/{id:s[0-9]+}/wfstatus/{jsonb}", a.PostSourceValue).Methods("POST")
	a.Router.HandleFunc("/source/{id:s[0-9]+}/{jsonb}", a.PostSourceJSON).Methods("POST")
	a.Router.HandleFunc("/source/{id:s[0-9]+}", a.GetSourceID).Methods("GET")
	a.Router.HandleFunc("/source/{id:s[0-9]+}", a.DeleteSourceID).Methods("DELETE")
	a.Router.HandleFunc("/source", a.GetSource).Methods("GET")
	a.Router.HandleFunc("/source/find", a.FindSource).Methods("GET")
	// Dgima
	a.Router.HandleFunc("/drim", a.GetFilesToDgima).Methods("GET")
	a.Router.HandleFunc("/cassette", a.GetCassetteFiles).Methods("GET")
	a.Router.HandleFunc("/drim/{id}", a.GetDgimaBySource).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.PostDgimaID).Methods("PUT")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}/wfstatus/{jsonb}", a.PostDgimaValue).Methods("POST")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}/{jsonb}", a.PostDgimaJSON).Methods("POST")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.GetDgimaID).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:[0-9]+}", a.GetDgimaByID).Methods("GET")
	a.Router.HandleFunc("/dgima/{id:d[0-9]+}", a.DeleteDgimaID).Methods("DELETE")
	a.Router.HandleFunc("/dgima", a.GetDgima).Methods("GET")
	a.Router.HandleFunc("/dgima/find", a.FindDgima).Methods("GET")
	a.Router.HandleFunc("/dgima/{jsonb}", a.FindDgimaByJSON).Methods("GET")
	// Trimmer
	a.Router.HandleFunc("/trim", a.GetFilesToTrim).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.PostTrimmerID).Methods("PUT")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}/wfstatus/{jsonb}", a.PostTrimmerValue).Methods("POST")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}/{jsonb}", a.PostTrimmerJSON).Methods("POST")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.GetTrimmerID).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:[0-9]+}", a.GetTrimmerByID).Methods("GET")
	a.Router.HandleFunc("/trimmer/{id:t[0-9]+}", a.DeleteTrimmerID).Methods("DELETE")
	a.Router.HandleFunc("/trimmer", a.GetTrimmer).Methods("GET")
	a.Router.HandleFunc("/trimmer/find", a.FindTrimmer).Methods("GET")
	a.Router.HandleFunc("/trimmer/{jsonb}", a.FindTrimmerByJSON).Methods("GET")
	// Aricha
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.PostArichaID).Methods("PUT")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}/wfstatus/{jsonb}", a.PostArichaValue).Methods("POST")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}/{jsonb}", a.PostArichaJSON).Methods("POST")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.GetArichaID).Methods("GET")
	a.Router.HandleFunc("/aricha/{id:[0-9]+}", a.GetArichaByID).Methods("GET")
	a.Router.HandleFunc("/aricha/{id:a[0-9]+}", a.DeleteArichaID).Methods("DELETE")
	a.Router.HandleFunc("/aricha", a.GetAricha).Methods("GET")
	a.Router.HandleFunc("/bdika", a.GetBdika).Methods("GET")
	a.Router.HandleFunc("/aricha/find", a.FindAricha).Methods("GET")
	a.Router.HandleFunc("/aricha/{jsonb}", a.FindArichaByJSON).Methods("GET")
	// Jobs
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.PostJobID).Methods("PUT")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}/wfstatus/{jsonb}", a.PostJobValue).Methods("POST")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}/{jsonb}", a.PostJobJSON).Methods("POST")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.GetJobID).Methods("GET")
	a.Router.HandleFunc("/jobs/{id:[0-9]+}", a.GetJobByID).Methods("GET")
	a.Router.HandleFunc("/jobs/{id:j[0-9]+}", a.DeleteJobID).Methods("DELETE")
	a.Router.HandleFunc("/jobs_list", a.GetListJobs).Methods("GET")
	a.Router.HandleFunc("/jobs", a.GetActiveJobs).Methods("GET")
	a.Router.HandleFunc("/jobs/find", a.FindJob).Methods("GET")
	a.Router.HandleFunc("/jobs/{jsonb}", a.FindJobByJSON).Methods("GET")
	// Products
	a.Router.HandleFunc("/products/{id:p[0-9]+}", a.PostProductID).Methods("PUT")
	a.Router.HandleFunc("/products/{id:p[0-9]+}/status", a.PostProductStatus).Methods("POST")
	a.Router.HandleFunc("/products/{id:p[0-9]+}/prop", a.PostProductProp).Methods("POST")
	a.Router.HandleFunc("/products/{id:p[0-9]+}/{prop}/{jsonb}", a.SetProductJSON).Methods("POST")
	//a.Router.HandleFunc("/products/{id:p[0-9]+}/{jsonb}", a.PostProductJSON).Methods("PUT")
	a.Router.HandleFunc("/products/{id:p[0-9]+}", a.GetProductID).Methods("GET")
	a.Router.HandleFunc("/products/{id:[0-9]+}", a.GetProductByID).Methods("GET")
	a.Router.HandleFunc("/products/{id:p[0-9]+}", a.DeleteProductID).Methods("DELETE")
	a.Router.HandleFunc("/products", a.GetListProducts).Methods("GET")
	//a.Router.HandleFunc("/products/{language}", a.GetActiveProducts).Methods("GET")
	a.Router.HandleFunc("/products/find", a.FindProduct).Methods("GET")
	//a.Router.HandleFunc("/products/{jsonb}", a.FindProductByJSON).Methods("GET")
	// Files
	a.Router.HandleFunc("/files/{id:f[0-9]+}", a.PostFile).Methods("PUT")
	a.Router.HandleFunc("/files/{id:f[0-9]+}", a.GetFile).Methods("GET")
	a.Router.HandleFunc("/files/{id:f[0-9]+}", a.DeleteFile).Methods("DELETE")
	a.Router.HandleFunc("/files/{id:f[0-9]+}/json/{jsonb}", a.PostFileJSON).Methods("PUT")
	a.Router.HandleFunc("/files/{id:f[0-9]+}/json/{jsonb}", a.PostFileValue).Methods("POST")
	a.Router.HandleFunc("/files/{id:f[0-9]+}/status/{jsonb}", a.PostFileStatus).Methods("POST")
	a.Router.HandleFunc("/files", a.GetFiles).Methods("GET")
	//a.Router.HandleFunc("/files/{language}", a.GetActiveFiles).Methods("GET")
	a.Router.HandleFunc("/files/find", a.FindFiles).Methods("GET")
	// Tasks
	a.Router.HandleFunc("/task", a.PostTask).Methods("POST")
	// Labels
	a.Router.HandleFunc("/labels", a.GetLabels).Methods("GET")
	a.Router.HandleFunc("/label/{id:[0-9]+}", a.GetLabel).Methods("GET")
	a.Router.HandleFunc("/labels/find", a.FindLabels).Methods("GET")
	// Cloud
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}", a.PostCloudID).Methods("PUT")
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}/status", a.PostCloudStatus).Methods("POST")
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}/prop", a.PostCloudProp).Methods("POST")
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}/{prop}/{jsonb}", a.SetCloudJSON).Methods("POST")
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}", a.GetCloudID).Methods("GET")
	a.Router.HandleFunc("/cloud/{id:[0-9]+}", a.GetCloudByID).Methods("GET")
	a.Router.HandleFunc("/cloud/{id:o[0-9]+}", a.DeleteCloudID).Methods("DELETE")
	a.Router.HandleFunc("/cloud", a.GetListClouds).Methods("GET")
	a.Router.HandleFunc("/cloud/find", a.FindCloud).Methods("GET")
	// State
	a.Router.HandleFunc("/states", a.GetStates).Methods("GET")
	a.Router.HandleFunc("/{tag}", a.GetStateByTag).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}", a.GetState).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.GetStateJSON).Methods("GET")
	a.Router.HandleFunc("/{tag}/{id}", a.PostState).Methods("PUT")
	a.Router.HandleFunc("/{tag}/{id}", a.UpdateState).Methods("POST")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.PostStateJSON).Methods("PUT")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.PostStateValue).Methods("POST")
	a.Router.HandleFunc("/{tag}/{id}", a.DeleteState).Methods("DELETE")
	a.Router.HandleFunc("/{tag}/{id}/{jsonb}", a.DeleteStateJSON).Methods("DELETE")
}

func (a *App) initMQTT() {
	if common.SERVER != "" {
		//a.InitLogMQTT()
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("ssl://%s", common.SERVER))
		opts.SetClientID("wfdb_mqtt_client")
		opts.SetUsername(common.USERNAME)
		opts.SetPassword(common.PASSWORD)
		opts.SetAutoReconnect(true)
		opts.SetOnConnectHandler(a.SubMQTT)
		opts.SetConnectionLostHandler(a.LostMQTT)
		a.Msg = mqtt.NewClient(opts)
		if token := a.Msg.Connect(); token.Wait() && token.Error() != nil {
			err := token.Error()
			log.Fatal().Str("source", "MQTT").Err(err).Msg("initialize mqtt listener")
		}
	}
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
