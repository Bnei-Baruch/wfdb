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

	"strconv"
)

type metus struct {
	MetusID		int			`json:"metus_id"`
	FileName	string  	`json:"filename"`
	UPath		string  	`json:"unix_path"`
	WPath		string  	`json:"windows_path"`
	Title		string  	`json:"title"`
	Sha1		string  	`json:"sha1"`
	Size 		float64		`json:"size"`
	Language	string		`json:"language"`
	Height		string		`json:"height"`
	Width		string		`json:"width"`
	Original	string		`json:"original"`
	Aspect		string		`json:"aspect"`
	Lecturer	string		`json:"lecturer"`
	Format		string		`json:"format"`
	Collection	string		`json:"collection"`
	Type		string		`json:"type"`
	Descr		string		`json:"desc"`
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

		err = c.getMetusMeta(db, c.MetusID, key)
		if err != nil {
			return nil, err
		}

		objects = append(objects, c)
	}

	return objects, nil
}

func (c *metus) getMetusMeta(db *sql.DB, MetusID int, key string) error {

	id := strconv.Itoa(MetusID)

	q := `SELECT (SELECT Value_String AS Height FROM dbo.METADATA_0 WHERE FieldID=1134 AND ObjectID=`+id+`) Height,
		(SELECT Value_String AS Collection FROM dbo.METADATA_0 WHERE FieldID=1000060 AND ObjectID=`+id+`) Collection,
		(SELECT Value_String AS Type FROM dbo.METADATA_0 WHERE FieldID=1000054 AND ObjectID=`+id+`) Type,
		(SELECT Value_String AS Descr FROM dbo.METADATA_0 WHERE FieldID=1000055 AND ObjectID=`+id+`) Descr,
		(SELECT Value_String AS Width FROM dbo.METADATA_0 WHERE FieldID=1133 AND ObjectID=`+id+`) Width,
		(SELECT Value_String AS Aspect FROM dbo.METADATA_0 WHERE FieldID=1082 AND ObjectID=`+id+`) Aspect,
		(SELECT Value_String AS Format FROM dbo.METADATA_0 WHERE FieldID=1142 AND ObjectID=`+id+`) Format,
		(SELECT Value_String AS Original FROM dbo.METADATA_0 WHERE FieldID=1000049 AND ObjectID=`+id+`) Original,
		(SELECT Value_String AS Lecturer FROM dbo.METADATA_0 WHERE FieldID=1000050 AND ObjectID=`+id+`) Lecturer,
		(SELECT Value_String AS FileName FROM dbo.METADATA_0 WHERE FieldID=2028 AND ObjectID=`+id+`) FileName,
		(SELECT Value_Number AS Size FROM dbo.METADATA_0 WHERE FieldID=1032 AND ObjectID=`+id+`) Size,
		(SELECT Value_String AS WPath FROM dbo.METADATA_0 WHERE FieldID=1033 AND ObjectID=`+id+`) WPath,
		(SELECT Value_String AS UPath FROM dbo.METADATA_0 WHERE FieldID=1034 AND ObjectID=`+id+`) UPath,
		(SELECT Value_String AS Title FROM dbo.METADATA_0 WHERE FieldID=1009 AND ObjectID=`+id+`) Title,
		(SELECT Value_String AS Language FROM dbo.METADATA_0 WHERE FieldID=1045 AND ObjectID=`+id+`) Language;`

	err := db.QueryRow(q).Scan(&c.Height, &c.Collection, &c.Type, &c.Descr, &c.Width, &c.Aspect, &c.Format, &c.Original, &c.Lecturer, &c.FileName, &c.Size, &c.WPath, &c.UPath, &c.Title, &c.Language)

	if err != nil {
		return err
	}

	c.WPath = strings.Replace(c.WPath, "\\\\", "\\", -1)
	c.UPath = strings.Replace(c.UPath, "\\", "/", -1)
	c.UPath = strings.Replace(c.UPath, ":", "-", -1)
	c.UPath = "/mnt/metus/" + strings.Replace(c.UPath, "/", "", 1)

	//err = c.getJSON("http://wfrp.bbdomain.org:8080/insert/find?key=file_name&value="+strings.TrimSuffix(c.FileName,path.Ext(c.FileName)))
	//if err != nil {
	//	return err
	//}

	if(key == "sha1") {

		f, err := os.Open(c.UPath)
		if err != nil {
			return err
		}

		h := sha1.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		c.Sha1 = hex.EncodeToString(h.Sum(nil))

		err = c.getJSON("http://wfrp.bbdomain.org:8080/insert/find?key=insert_name&value="+c.FileName)
		if err != nil {
			return err
		}
	}

	return nil
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