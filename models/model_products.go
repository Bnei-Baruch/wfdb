package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Products struct {
	ID          int         `json:"id"`
	ProductID   string      `json:"product_id"`
	Date        string      `json:"date"`
	Language    string      `json:"language"`
	OriginLang  string      `json:"original_language"`
	Pattern     string      `json:"pattern"`
	ProductName string      `json:"product_name"`
	ProductType string      `json:"product_type"`
	Parent      interface{} `json:"parent"`
	Line        interface{} `json:"line"`
	Wfstatus    interface{} `json:"wfstatus"`
}

func FindProduct(db *sql.DB, key string, value string) ([]Products, error) {
	sqlStatement := `SELECT * FROM products WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY product_id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var parent, line, wfstatus []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func FindProductByJSON(db *sql.DB, ep string, key string, value string) ([]Products, error) {

	sqlStatement := fmt.Sprintf("SELECT id, product_id, date, language, original_language, pattern, product_name, product_type, parent, line, wfstatus FROM products WHERE %s ->> '%s' = '%s' ORDER BY product_id;", ep, key, value)
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var parent, line, wfstatus []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}

func GetListProducts(db *sql.DB, start, count int) ([]Products, error) {
	rows, err := db.Query(
		"SELECT id, product_id, date, language, original_language, pattern, product_name, product_type, parent, line, wfstatus FROM products ORDER BY product_id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var parent, line, wfstatus []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func GetActiveProducts(db *sql.DB, language string) ([]Products, error) {
	rows, err := db.Query(
		"SELECT * FROM products WHERE wfstatus ->> 'removed' = 'false' AND line ->> 'language' = $1 ORDER BY product_id", language)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	objects := []Products{}

	for rows.Next() {
		var t Products
		var parent, line, wfstatus []byte
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus); err != nil {
			return nil, err
		}
		json.Unmarshal(parent, &t.Parent)
		json.Unmarshal(line, &t.Line)
		json.Unmarshal(wfstatus, &t.Wfstatus)
		objects = append(objects, t)
	}

	return objects, nil
}

func (t *Products) GetProductID(db *sql.DB) error {
	var parent, line, wfstatus []byte

	err := db.QueryRow("SELECT id, product_id, date, language, original_language, pattern, product_name, product_type, parent, line, wfstatus FROM products WHERE product_id = $1",
		t.ProductID).Scan(
		&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Products) GetProductByID(db *sql.DB) error {
	var parent, line, wfstatus []byte

	err := db.QueryRow("SELECT id, product_id, date, pattern, product_name, product_type, parent, line, departments, proxy, product, wfstatus FROM products WHERE id = $1",
		t.ID).Scan(&t.ID, &t.ProductID, &t.Date, &t.Language, &t.OriginLang, &t.Pattern, &t.ProductName, &t.ProductType, &parent, &line, &wfstatus)

	json.Unmarshal(parent, &t.Parent)
	json.Unmarshal(line, &t.Line)
	json.Unmarshal(wfstatus, &t.Wfstatus)
	if err != nil {
		return err
	}

	return err
}

func (t *Products) PostProductID(db *sql.DB) error {
	parent, _ := json.Marshal(t.Parent)
	line, _ := json.Marshal(t.Line)
	wfstatus, _ := json.Marshal(t.Wfstatus)

	err := db.QueryRow(
		"INSERT INTO products(product_id, date, language, original_language, pattern, product_name, product_type, parent, line, wfstatus) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (product_id) DO UPDATE SET (product_id, date, language, original_language, pattern, product_name, product_type, parent, line, wfstatus) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) WHERE products.product_id = $1 RETURNING id",
		t.ProductID, t.Date, t.Language, t.OriginLang, t.Pattern, t.ProductName, t.ProductType, parent, line, wfstatus).Scan(&t.ID)

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

	_, err := db.Exec("UPDATE products SET wfstatus = wfstatus || json_build_object($3::text, $2::bool)::jsonb WHERE product_id=$1",
		t.ProductID, value, key)

	return err
}

func (t *Products) DeleteProductID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE product_id=$1", t.ProductID)

	return err
}
