package cmd

import (
	"os"

	"github.com/Bnei-Baruch/wfdb/api"
)

func Init() {
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	accountsUrl := os.Getenv("ACC_URL")
	skipAuth := os.Getenv("SKIP_AUTH") == "true"

	a := api.App{}
	a.Initialize(accountsUrl, skipAuth)
	a.InitDB(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("METUS_DB_HOST"),
		os.Getenv("METUS_DB_USERNAME"),
		os.Getenv("METUS_DB_PASSWORD"),
		os.Getenv("METUS_DB_NAME"))
	a.Run(listenAddress)
}
