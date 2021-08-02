package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Products struct {
	ID          int         `json:"id"`
	ProductID   string      `json:"product_id"`
	Date        string      `json:"date"`
	Language    string      `json:"language"`
	Pattern     string      `json:"pattern"`
	TypeID      string      `json:"type_id"`
	ProductName string      `json:"product_name"`
	ProductType string      `json:"product_type"`
	I18n        interface{} `json:"i18n"`
	Parent      interface{} `json:"parent"`
	Line        interface{} `json:"line"`
	Props       interface{} `json:"properties"`
	FilmDate    string      `json:"film_date"`
}

func FindProduct(db *sql.DB, values url.Values) ([]Products, error) {

	var where []string
	sqlStatement := `SELECT * FROM products WHERE properties ->> 'removed' = 'false'`
	collection_uid := "0"
	limit := "10"
	offset := "0"

	for k, v := range values {
		if k == "limit" {
			limit = v[0]
			continue
		}
		if k == "offset" {
			offset = v[0]
			continue
		}
		if k == "collection_uid" {
			collection_uid = v[0]
			continue
		}
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, k, v[0]))
	}

	if len(where) == 0 && collection_uid != "0" {
		sqlStatement = sqlStatement + ` AND line ->> 'collection_uid' = '` + collection_uid + `' AND `
	} else if len(where) > 0 && collection_uid != "0" {
		sqlStatement = sqlStatement + ` AND line ->> 'collection_uid' = '` + collection_uid + `' AND ` + strings.Join(where, " AND ")
	} else if len(where) > 0 && collection_uid == "0" {
		sqlStatement = sqlStatement + ` AND ` + strings.Join(where, " AND ")
	}

	sqlStatement = sqlStatement + fmt.Sprintf(` ORDER BY product_id DESC LIMIT %s OFFSET %s`, limit, offset)

	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var i18n, parent, line, properties []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate); err != nil {
			return nil, err
		}
		json.Unmarshal(i18n, &t.I18n)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindProductByJSON(db *sql.DB, ep string, key string, value string) ([]Products, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM products WHERE %s ->> '%s' = '%s' ORDER BY product_id;", ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var i18n, parent, line, properties []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate); err != nil {
			return nil, err
		}
		json.Unmarshal(i18n, &t.I18n)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetListProducts(db *sql.DB, start, count int) ([]Products, error) {
	rows, err := db.Query(
		"SELECT * FROM products ORDER BY product_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var i18n, parent, line, properties []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate); err != nil {
			return nil, err
		}
		json.Unmarshal(i18n, &t.I18n)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func GetActiveProducts(db *sql.DB, language string) ([]Products, error) {
	rows, err := db.Query(
		"SELECT * FROM products WHERE properties ->> 'removed' = 'false' AND line ->> $1 IS NOT NULL ORDER BY product_id", language)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var i18n, parent, line, properties []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate); err != nil {
			return nil, err
		}
		json.Unmarshal(i18n, &t.I18n)
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(properties, &t.Props)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Products) GetProductID(db *sql.DB) error {
	var i18n, parent, line, properties []byte

	err := db.QueryRow("SELECT * FROM products WHERE product_id = $1",
		t.ProductID).Scan(
		&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate)
	json.Unmarshal(i18n, &t.I18n)
	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Products) GetProductByID(db *sql.DB) error {
	var i18n, parent, line, properties []byte

	err := db.QueryRow("SELECT * FROM products WHERE id = $1",
		t.ID).Scan(&t.ID, &t.ProductID, &t.Date, &t.Language, &t.TypeID, &t.Pattern, &t.ProductName, &t.ProductType, &i18n, &parent, &line, &properties, &t.FilmDate)
	json.Unmarshal(i18n, &t.I18n)
	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(properties, &t.Props)
	if err != nil {
		return err
	}

	return err
}

func (t *Products) PostProductID(db *sql.DB) error {
	i18n, _ := json.Marshal(t.I18n)
	parent, _ := json.Marshal(t.Parent)
	line, _ := json.Marshal(t.Line)
	properties, _ := json.Marshal(t.Props)

	err := db.QueryRow(
		"INSERT INTO products(product_id, date, language, type_id, pattern, product_name, product_type, i18n, parent, line, properties, film_date) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (product_id) DO UPDATE SET (product_id, date, language, type_id, pattern, product_name, product_type, i18n, parent, line, properties, film_date) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) WHERE products.product_id = $1 RETURNING id",
		t.ProductID, t.Date, t.Language, t.TypeID, t.Pattern, t.ProductName, t.ProductType, i18n, parent, line, properties, t.FilmDate).Scan(&t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (t *Products) PostProductJSON(db *sql.DB, jsonb interface{}, key string) error {
	v, _ := json.Marshal(jsonb)

	sqlStatement := `UPDATE products SET ` + key + ` = $2 WHERE product_id=$1;`
	_, err := db.Exec(sqlStatement, t.ProductID, v)

	return err
}

func (t *Products) SetProductJSON(db *sql.DB, value interface{}, key string, prop string) error {

	v, _ := json.Marshal(value)
	sqlCmd := "UPDATE products SET " + prop + " = " + prop + " || json_build_object($3::text, $2::jsonb)::jsonb WHERE product_id=$1"
	_, err := db.Exec(sqlCmd, t.ProductID, v, key)

	return err
}

func (t *Products) PostProductStatus(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE products SET properties = properties || json_build_object($3::text, $2::bool)::jsonb WHERE product_id=$1",
		t.ProductID, value, key)

	return err
}

func (t *Products) PostProductProp(db *sql.DB, value, key string) error {

	_, err := db.Exec("UPDATE products SET properties = properties || json_build_object($3::text, $2::text)::jsonb WHERE product_id=$1",
		t.ProductID, value, key)

	return err
}

func (t *Products) DeleteProductID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE product_id=$1", t.ProductID)

	return err
}
