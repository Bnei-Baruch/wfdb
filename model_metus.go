// model_metus.go

package main

import (
	"database/sql"
	"strings"
	"os"
	"crypto/sha1"
	"io"
	"encoding/hex"
)

type metus struct {
	MetusID    int		`json:"metus_id"`
	FileName  string  	`json:"filename"`
	UPath  string  		`json:"unix_path"`
	WPath  string  		`json:"windows_path"`
	Title  string  		`json:"title"`
	Sha1  string  		`json:"sha1"`
}

func findMetus(db *sql.DB, key string, value string) ([]metus, error) {
	sqlStatement := `SELECT DISTINCT ObjectID FROM METADATA_0 WHERE Value_String LIKE '%`+value+`%' AND FieldID=2028`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []metus{}

	for rows.Next() {
		var c metus
		if err := rows.Scan(&c.MetusID); err != nil {
			return nil, err
		}
		//fmt.Println("  Select db:", c.ID)

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=2028", c.MetusID).Scan(&c.FileName)
		if err != nil {
			return nil, err
		}
		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1033", c.MetusID).Scan(&c.WPath)
		if err != nil {
			return nil, err
		}
		c.WPath = strings.Replace(c.WPath, "\\\\", "\\", -1)

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1034", c.MetusID).Scan(&c.UPath)
		if err != nil {
			return nil, err
		}
		c.UPath = strings.Replace(c.UPath, "\\", "/", -1)
		c.UPath = strings.Replace(c.UPath, ":", "-", -1)
		c.UPath = "/mnt/metus/" + strings.Replace(c.UPath, "/", "", 1)

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1009", c.MetusID).Scan(&c.Title)
		if err != nil {
			return nil, err
		}

		if(key == "sha1") {

			f, err := os.Open(c.UPath)
			if err != nil {
				return nil, err
			}

			h := sha1.New()
			if _, err := io.Copy(h, f); err != nil {
				return nil, err
			}

			c.Sha1 = hex.EncodeToString(h.Sum(nil))
		}
		objects = append(objects, c)
	}

	return objects, nil
}