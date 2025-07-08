package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
)

func InsertEntry(entry models.Entry) error {
	db := variables.DB
	_, err := db.Exec("INSERT INTO entries (record_id, data, created_at) VALUES ($1, $2, $3)",
		entry.Record_id, entry.Data, entry.Created_at)
	if err != nil {
		log.Printf("Ошибка при вставке записи: %v", err)
	} else {
		log.Printf("Record вставлен: %+v", entry)
	}
	return err
}

func SelectEntry(limit, offset int) ([]models.Entry, error) {
	var entries []models.Entry
	db := variables.DB
	rows, err := db.Query("SELECT id, record_id, data, created_at FROM entries ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Printf("DB query error: %v", err)
		return entries, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(&entry.Id, &entry.Record_id, &entry.Data, &entry.Created_at); err != nil {
			log.Printf("Ошибка сканирования записи: %v", err)
			continue
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка после итерации по rows: %v", err)
		return entries, err
	}
	return entries, nil
}

func SelectEntryByRecordId(id int) (models.Entry, error) {
	var entry models.Entry
	db := variables.DB
	row := db.QueryRow("Select id , Record_id, Data , Created_at FROM entries where Record_id = $1", id)
	err := row.Scan(&entry.Id, &entry.Record_id, &entry.Data, &entry.Created_at)
	if err != nil {
		log.Printf("Ошибка при получении записи: %v", err)
	}
	return entry, err
}
func DeleteEntryById(entry_id int) error {
	db := variables.DB
	_, err := db.Exec("DELETE FROM entries WHERE id = $1", entry_id)
	if err != nil {
		log.Printf(err.Error())
	}
	if err == nil {
		log.Printf("Запись успешно удалена: ID = %v", entry_id)
	}
	return err
}
