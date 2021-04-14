package models

import (
	"database/sql"
)

type Archive struct {
	ID        int    `json:"id"`
	ArchiveID string `json:"archive_id"`
	Date      string `json:"date"`
	FileName  string `json:"file_name"`
	Language  string `json:"language"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
	Sha1      string `json:"sha1"`
	SendID    string `json:"send_id"`
	Source    string `json:"source"`
}

func FindArFiles(db *sql.DB, key string, value string) ([]Archive, error) {
	sqlStatement := `SELECT id, archive_id, date, file_name, language, extension, size, sha1, send_id, source FROM archive WHERE ` + key + ` LIKE '%` + value + `%' ORDER BY id`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	arfiles := []Archive{}

	for rows.Next() {
		var a Archive
		if err := rows.Scan(&a.ID, &a.ArchiveID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.SendID, &a.Source); err != nil {
			return nil, err
		}
		arfiles = append(arfiles, a)
	}

	return arfiles, nil
}

func GetArFiles(db *sql.DB, start, count int) ([]Archive, error) {
	rows, err := db.Query(
		"SELECT id, archive_id, date, file_name, language, extension, size, sha1, send_id, source FROM archive ORDER BY id DESC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	arfiles := []Archive{}

	for rows.Next() {
		var a Archive
		if err := rows.Scan(&a.ID, &a.ArchiveID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.SendID, &a.Source); err != nil {
			return nil, err
		}
		arfiles = append(arfiles, a)
	}

	return arfiles, nil
}

func (a *Archive) GetArFile(db *sql.DB) error {

	return db.QueryRow("SELECT id, archive_id, date, file_name, language, extension, size, sha1, send_id, source FROM archive WHERE archive_id = $1",
		a.ArchiveID).Scan(&a.ID, &a.ArchiveID, &a.Date, &a.FileName, &a.Language, &a.Extension, &a.Size, &a.Sha1, &a.SendID, &a.Source)
}

func (a *Archive) PostArFile(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO archive(archive_id, date, file_name, language, extension, size, sha1, send_id, source) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (archive_id) DO UPDATE SET (archive_id, date, file_name, language, extension, size, sha1, send_id, source) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE archive.archive_id = $1 RETURNING id",
		a.ArchiveID, a.Date, a.FileName, a.Language, a.Extension, a.Size, a.Sha1, a.SendID, a.Source).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (a *Archive) UpdateArFile(db *sql.DB) error {

	_, err :=
		db.Exec("UPDATE archive SET (archive_id, date, file_name, language, extension, size, sha1, send_id, source) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE aricha_id=$1",
			a.ArchiveID, a.Date, a.FileName, a.Language, a.Extension, a.Size, a.Sha1, a.SendID, a.Source)

	return err
}

func (a *Archive) DeleteArFile(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM archive WHERE archive_id=$1", a.ArchiveID)

	return err
}
