// model_metus.go

package main

import (
	"database/sql"
	"strings"
	"os"
	"crypto/sha1"
	"io"
	"encoding/hex"
	"net/http"
	"fmt"
	"encoding/json"
)

type metus struct {
	MetusID		int			`json:"metus_id"`
	FileName	string  	`json:"filename"`
	UPath		string  	`json:"unix_path"`
	WPath		string  	`json:"windows_path"`
	Title		string  	`json:"title"`
	Sha1		string  	`json:"sha1"`
	Language	string		`json:"language"`
	Height		string		`json:"height"`
	Width		string		`json:"width"`
	Original	string		`json:"original"`
	Aspect		string		`json:"aspect"`
	Lecturer	string		`json:"lecturer"`
	Format		string		`json:"format"`
	Collection	string		`json:"collection"`
	Type		string		`json:"type"`
	Desc		string		`json:"desc"`
	Workflow []interface{} `json:"workflow"`
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
		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1134", c.MetusID).Scan(&c.Height)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1000060", c.MetusID).Scan(&c.Collection)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1000054", c.MetusID).Scan(&c.Type)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1000055", c.MetusID).Scan(&c.Desc)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1133", c.MetusID).Scan(&c.Width)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1082", c.MetusID).Scan(&c.Aspect)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1142", c.MetusID).Scan(&c.Format)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1000049", c.MetusID).Scan(&c.Original)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1000050", c.MetusID).Scan(&c.Lecturer)
		if err != nil {
			return nil, err
		}

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1133", c.MetusID).Scan(&c.Width)
		if err != nil {
			return nil, err
		}

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

		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1045", c.MetusID).Scan(&c.Language)
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

		//err = c.getJSON("http://wfrp.bbdomain.org:8080/aricha/find?key=file_name&value="+strings.TrimSuffix(c.FileName,path.Ext(c.FileName)))
		//if err != nil {
		//	return nil, err
		//}

		objects = append(objects, c)
	}

	return objects, nil
}

func (c *metus) getJSON(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}

	if err != nil {
		return fmt.Errorf("cannot read JSON: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&c.Workflow)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %v", err)
	}

	return nil
}