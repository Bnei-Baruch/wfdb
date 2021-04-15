package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"strconv"
)

type Metus struct {
	MetusID    int           `json:"metus_id"`
	FileName   string        `json:"filename"`
	UPath      string        `json:"unix_path"`
	WPath      string        `json:"windows_path"`
	Title      string        `json:"title"`
	Sha1       string        `json:"sha1"`
	Size       float64       `json:"size"`
	Language   string        `json:"language"`
	Original   string        `json:"original"`
	Lecturer   string        `json:"lecturer"`
	Collection string        `json:"collection"`
	Type       string        `json:"type"`
	Descr      string        `json:"desc"`
	Workflow   []interface{} `json:"workflow"`
}

func FindMetus(db *sql.DB, key string, value string) ([]Metus, error) {
	sqlStatement := `SELECT DISTINCT ObjectID FROM METADATA_0 WHERE Value_String LIKE '%` + value + `%' AND FieldID=2028`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Metus{}

	for rows.Next() {
		var c Metus
		if err := rows.Scan(&c.MetusID); err != nil {
			return nil, err
		}

		err = c.GetMetusMeta(db, c.MetusID, key)
		if err != nil {
			return nil, err
		}

		objects = append(objects, c)
	}

	return objects, nil
}

func (c *Metus) GetMetusMeta(db *sql.DB, MetusID int, key string) error {

	id := strconv.Itoa(MetusID)

	q := `SELECT 
		(SELECT Value_String AS Collection FROM dbo.METADATA_0 WHERE FieldID=1000060 AND ObjectID=` + id + `) Collection,
		(SELECT Value_String AS Type FROM dbo.METADATA_0 WHERE FieldID=1000054 AND ObjectID=` + id + `) Type,
		(SELECT Value_String AS Descr FROM dbo.METADATA_0 WHERE FieldID=1000055 AND ObjectID=` + id + `) Descr,
		(SELECT Value_String AS Original FROM dbo.METADATA_0 WHERE FieldID=1000049 AND ObjectID=` + id + `) Original,
		(SELECT Value_String AS Lecturer FROM dbo.METADATA_0 WHERE FieldID=1000050 AND ObjectID=` + id + `) Lecturer,
		(SELECT Value_String AS FileName FROM dbo.METADATA_0 WHERE FieldID=2028 AND ObjectID=` + id + `) FileName,
		(SELECT Value_Number AS Size FROM dbo.METADATA_0 WHERE FieldID=1032 AND ObjectID=` + id + `) Size,
		(SELECT Value_String AS WPath FROM dbo.METADATA_0 WHERE FieldID=1033 AND ObjectID=` + id + `) WPath,
		(SELECT Value_String AS UPath FROM dbo.METADATA_0 WHERE FieldID=1034 AND ObjectID=` + id + `) UPath,
		(SELECT Value_String AS Title FROM dbo.METADATA_0 WHERE FieldID=1009 AND ObjectID=` + id + `) Title,
		(SELECT Value_String AS Language FROM dbo.METADATA_0 WHERE FieldID=1045 AND ObjectID=` + id + `) Language;`

	err := db.QueryRow(q).Scan(&c.Collection, &c.Type, &c.Descr, &c.Original, &c.Lecturer, &c.FileName, &c.Size, &c.WPath, &c.UPath, &c.Title, &c.Language)

	if err != nil {
		return err
	}

	c.WPath = strings.Replace(c.WPath, "\\\\", "\\", -1)
	c.UPath = strings.Replace(c.UPath, "\\", "/", -1)
	c.UPath = strings.Replace(c.UPath, ":", "-", -1)
	c.UPath = "/mnt/metus/" + strings.Replace(c.UPath, "/", "", 1)

	//err = c.GetJSON("http://wfrp.bbdomain.org:8080/insert/find?key=file_name&value="+strings.TrimSuffix(c.FileName,path.Ext(c.FileName)))
	//if err != nil {
	//	return err
	//}

	if key == "sha1" {

		f, err := os.Open(c.UPath)
		if err != nil {
			return err
		}

		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		c.Sha1 = hex.EncodeToString(h.Sum(nil))

		err = c.GetJSON("http://wfrp.bbdomain.org:8080/insert/find?key=insert_name&value=" + c.FileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Metus) GetJSON(url string) error {

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
