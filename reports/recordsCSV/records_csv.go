package recordsCSV

import (
	db_ "awesomeProject1/internal/db"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func GenerateCSVRecords(w io.Writer, limit, offset int) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write([]string{"RecordId", "Timeout", "CreatedAt", "Status"})
	rows := [][]string{}
	records, err := db_.SelectRecord(limit, offset)
	if err != nil {
		log.Printf("Ошибка при генерации CSV: %v", err)
		return err
	}
	for _, record := range records {
		row := []string{}
		recordID := strconv.Itoa(int(record.RecordId))
		timeout := strconv.Itoa(int(record.Timeout))
		createdAt := record.CreatedAt.String()
		status := record.Status
		row = append(row, recordID)
		row = append(row, timeout)
		row = append(row, createdAt)
		row = append(row, status)
		rows = append(rows, row)
	}
	if err := writer.WriteAll(rows); err != nil {
		log.Printf("Ошибка записи CSV: %v", err)
		return err
	}
	return nil
}
