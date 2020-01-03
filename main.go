// main.go

package main

import "os"

func main() {
	a := App{}
	a.initOidc(os.Getenv("ACC_URL"))
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("METUS_DB_HOST"),
		os.Getenv("METUS_DB_USERNAME"),
		os.Getenv("METUS_DB_PASSWORD"),
		os.Getenv("METUS_DB_NAME"))
	a.Run(":8080")
}
