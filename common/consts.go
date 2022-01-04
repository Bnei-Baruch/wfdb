package common

import "os"

var (
	LISTEN_ADDRESS = os.Getenv("LISTEN_ADDRESS")
	ACC_URL        = os.Getenv("ACC_URL")
	SKIP_AUTH      = os.Getenv("SKIP_AUTH") == "true"

	APP_DB_USERNAME   = os.Getenv("APP_DB_USERNAME")
	APP_DB_PASSWORD   = os.Getenv("APP_DB_PASSWORD")
	APP_DB_NAME       = os.Getenv("APP_DB_NAME")
	APP_DB_HOST       = os.Getenv("APP_DB_HOST")
	APP_DB_PORT       = os.Getenv("APP_DB_PORT")
	METUS_DB_HOST     = os.Getenv("METUS_DB_HOST")
	METUS_DB_USERNAME = os.Getenv("METUS_DB_USERNAME")
	METUS_DB_PASSWORD = os.Getenv("METUS_DB_PASSWORD")
	METUS_DB_NAME     = os.Getenv("METUS_DB_NAME")

	SERVER   = os.Getenv("MQTT_URL")
	USERNAME = os.Getenv("MQTT_USER")
	PASSWORD = os.Getenv("MQTT_PASS")

	LogPath = os.Getenv("LOG_PATH")

	ServiceTopic = "wfdb/service/#"
)

const (
	ExtPrefix           = "kli/"
	ServiceDataTopic    = "wfdb/service/data/"
	MonitorIngestTopic  = "wfdb/service/ingest/monitor"
	MonitorTrimmerTopic = "wfdb/service/trimmer/monitor"
	MonitorArchiveTopic = "wfdb/service/archive/monitor"
	StateTrimmerTopic   = "wfdb/service/trimmer/state"
	StateDgimaTopic     = "wfdb/service/dgima/state"
	StateArichaTopic    = "wfdb/service/aricha/state"
	StateJobsTopic      = "wfdb/service/jobs/state"
	StateLangcheckTopic = "wfdb/service/langcheck/state"
)
