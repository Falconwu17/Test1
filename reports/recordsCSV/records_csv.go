package recordsCSV

import (
	db_ "awesomeProject1/internal/db"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func GenerateCSVRecords(w io.Writer) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write([]string{"Record_id", "Timeout", "Handler_type", "Created_at", "Status"})
	rows := [][]string{}
	records, err := db_.SelectRecord(100, 0)
	if err != nil {
		log.Printf("Ошибка при генерации CSV: %v", err)
		return
	}
	for _, record := range records {
		row := []string{}
		recordID := strconv.Itoa(int(record.Record_id))
		timeout := strconv.Itoa(int(record.Timeout))
		handlerType := record.Handler_type
		createdAt := record.Created_at.String()
		status := record.Status
		row = append(row, recordID)
		row = append(row, timeout)
		row = append(row, handlerType)
		row = append(row, createdAt)
		row = append(row, status)
		rows = append(rows, row)
	}
	if err := writer.WriteAll(rows); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
