package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
)

func InsertRecord(record models.Record) error {
	db := variables.DB
	_, err := db.Exec("INSERT INTO records ( Timeout, Created_at ,Status) VALUES ($1, $2, $3)",
		record.Timeout, record.Created_at, record.Status)
	if err != nil {
		log.Printf("Ошибка при вставке записи: %v", err)
	} else {
		log.Printf("Record вставлен: %+v", record)
	}
	return err
}

func CheckRecordExists(id int) (bool, error) {
	var exists bool
	err := variables.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM records WHERE record_id = $1)", id).Scan(&exists)
	if err != nil {
		log.Printf("Ошибка при проверке существования записи: %v", err)
		return false, err
	}
	return exists, nil
}

func SelectRecord(limit, offset int) ([]models.Record, error) {
	records := []models.Record{}
	db := variables.DB
	query := "SELECT record_id, timeout, created_at, status FROM records ORDER BY created_at desc limit $1 offset $2"
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Printf("DB query error: %v", err)
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.Record
		if err := rows.Scan(&record.Record_id, &record.Timeout, &record.Created_at, &record.Status); err != nil {
			log.Printf("Ошибка сканирования записи: %v", err)
			continue
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Ошибка после итерации по rows: %v", err)
		return records, err
	}
	return records, nil
}
func SelectRecordById(record_id int) (models.Record, error) {
	db := variables.DB
	var record models.Record
	rows := db.QueryRow("SELECT record_id, timeout,  created_at, status FROM records WHERE Record_id = $1", record_id)
	err := rows.Scan(&record.Record_id, &record.Timeout, &record.Created_at, &record.Status)
	if err != nil {
		log.Printf("Ошибка при получении записи: %v", err)
	}
	return record, err
}
func DeleteRecordById(record_id int) error {
	db := variables.DB
	_, err := db.Exec("DELETE FROM records WHERE Record_id = $1", record_id)
	if err != nil {
		log.Printf(err.Error())
	}
	if err == nil {
		log.Printf("Запись успешно удалена: ID = %v", record_id)
	}
	return err
}
