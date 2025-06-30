package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
)

func InsertRecord(record models.Record) {
	db := variables.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("INSERT INTO records ( Timeout, Handler_type ,Created_at ,Status) VALUES ($1, $2, $3, $4)",
		record.Timeout, record.Handler_type, record.Created_at, record.Status)
	if err != nil {
		log.Printf(err.Error())
	}
}

func SelectRecord() []models.Record {
	records := []models.Record{}
	db := variables.ConnectToDB()
	defer db.Close()
	rows, err := db.Query("SELECT record_id, timeout, handler_type, created_at, status FROM records")
	if err != nil {
		log.Printf(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Record
		err = rows.Scan(&record.Record_id, &record.Timeout, &record.Handler_type, &record.Created_at, &record.Status)
		if err != nil {
			log.Printf(err.Error())
		}
		records = append(records, record)
	}
	return records
}
func SelectRecordById(record_id int) (models.Record, error) {
	db := variables.ConnectToDB()
	defer db.Close()
	var record models.Record
	rows := db.QueryRow("SELECT record_id, timeout, handler_type, created_at, status FROM records WHERE Record_id = $1", record_id)
	err := rows.Scan(&record.Record_id, &record.Timeout, &record.Handler_type, &record.Created_at, &record.Status)
	if err != nil {
		log.Printf("Ошибка при получении записи: %v", err)
	}
	return record, err
}
func DeleteRecordById(record_id int) error {
	db := variables.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("DELETE FROM records WHERE Record_id = $1", record_id)
	if err != nil {
		log.Printf(err.Error())
	}
	if err == nil {
		log.Printf("Record by ID deleted susheful: %v", record_id)
	}
	return err
}
