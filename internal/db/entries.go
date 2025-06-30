package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
)

func InsertEntry(entry models.Entry) {
	db := variables.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("INSERT INTO entries (record_id, data, created_at) VALUES ($1, $2, $3)",
		entry.Record_id, entry.Data, entry.Created_at)
	if err != nil {
		log.Printf(err.Error())
	}
}

func SelectEntry() []models.Entry {
	var entries []models.Entry
	db := variables.ConnectToDB()
	defer db.Close()
	rows, err := db.Query("SELECT id ,Record_id, Data , Created_at FROM entries")
	if err != nil {
		log.Printf(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var entry models.Entry
		err := rows.Scan(&entry.Id, &entry.Record_id, &entry.Data, &entry.Created_at)
		if err != nil {
			log.Printf(err.Error())
		}
		entries = append(entries, entry)
	}
	return entries
}
func SelectEntryByRecordId(id int) (models.Entry, error) {
	var entry models.Entry
	db := variables.ConnectToDB()
	defer db.Close()
	row := db.QueryRow("Select id , Record_id, Data , Created_at FROM entries where Record_id = $1", id)
	err := row.Scan(&entry.Id, &entry.Record_id, &entry.Data, &entry.Created_at)
	if err != nil {
		log.Printf(err.Error())
	}
	return entry, err
}
func DeleteEntryById(entry_id int) error {
	db := variables.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("DELETE FROM entries WHERE id = $1", entry_id)
	if err != nil {
		log.Printf(err.Error())
	}
	return err
}
