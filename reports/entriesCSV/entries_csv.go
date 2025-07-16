package entriesCSV

import (
	db_ "awesomeProject1/internal/db"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func GenerateCSVEntry(w io.Writer, limit, offset int) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write([]string{"Id", "RecordId", "Data", "CreatedAt"})
	entries, err := db_.SelectEntry(limit, offset)
	if err != nil {
		log.Printf("Ошибка при получении entries: %v", err)
		return
	}
	rows := [][]string{}
	for _, entry := range entries {
		row := []string{}
		id := strconv.Itoa(int(entry.Id))
		recordID := strconv.Itoa(int(entry.RecordId))
		createdAT := entry.CreatedAt.Format("2006-01-02 15:04:05")
		data := string(entry.Data)
		row = append(row, id)
		row = append(row, recordID)
		row = append(row, data)
		row = append(row, createdAT)
		rows = append(rows, row)
	}
	if err := writer.WriteAll(rows); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	return
}
