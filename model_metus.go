// model_metus.go

package main

import (
	"database/sql"
)

type metus struct {
	MetusID    int		`json:"metus_id"`
	FileName  string  	`json:"filename"`
	UPath  string  		`json:"unix_path"`
	WPath  string  		`json:"windows_path"`
	Title  string  		`json:"title"`
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
		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1034", c.MetusID).Scan(&c.UPath)
		if err != nil {
			return nil, err
		}
		err = db.QueryRow("SELECT Value_String FROM dbo.METADATA_0 WHERE ObjectID=$1 AND FieldID=1009", c.MetusID).Scan(&c.Title)
		if err != nil {
			return nil, err
		}
		objects = append(objects, c)
	}

	return objects, nil
}