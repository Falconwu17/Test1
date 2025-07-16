package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"fmt"
	"log"
	"time"
)

func InsertRecord(record *models.Record) error {
	db := variables.DB
	err := db.QueryRow(
		"INSERT INTO records (timeout, created_at, status) VALUES ($1, $2, $3) RETURNING record_id",
		record.Timeout, record.CreatedAt, record.Status,
	).Scan(&record.RecordId)

	if err != nil {
		log.Printf("Ошибка при вставке записи: %v", err)
		return err
	}
	if record.RecordId != 0 {
		log.Printf("Попытка задать ID вручную: %+v", *record)
	}
	log.Printf("Record вставлен: %+v", *record)
	return nil
}
func ensureDefaultRecord() error {
	exists, err := CheckRecordExists(1)
	if err != nil {
		log.Printf("<UNK> <UNK> <UNK> <UNK>: %v", err)
	}
	if !exists {
		r := models.Record{
			Timeout:   60,
			CreatedAt: time.Now(),
			Status:    "now",
		}
		r.SetDefault()
		err := InsertRecord(&r)
		if err != nil {
			return err
		}
		log.Println("Создан record_id=1 по умолчанию:")
	}
	return nil
}
func CheckRecord() error {
	if err := ensureDefaultRecord(); err != nil {
		fmt.Println("Ошибка при создании record 1:", err)
	}
	return ensureDefaultRecord()
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
		if err := rows.Scan(&record.RecordId, &record.Timeout, &record.CreatedAt, &record.Status); err != nil {
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
func SelectRecordById(recordId int) (models.Record, error) {
	db := variables.DB
	var record models.Record
	rows := db.QueryRow("SELECT record_id, timeout,  created_at, Status FROM records WHERE record_id= $1", recordId)
	err := rows.Scan(&record.RecordId, &record.Timeout, &record.CreatedAt, &record.Status)
	if err != nil {
		log.Printf("Ошибка при получении записи: %v", err)
	}
	return record, err
}
func DeleteRecordById(recordId int) error {
	db := variables.DB
	_, err := db.Exec("DELETE FROM records WHERE record_id = $1", recordId)
	if err != nil {
		log.Printf(err.Error())
	}
	if err == nil {
		log.Printf("Запись успешно удалена: ID = %v", recordId)
	}
	return err
}
