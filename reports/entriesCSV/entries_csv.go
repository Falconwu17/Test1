package entriesCSV

import (
	db_ "awesomeProject1/internal/db"
	"encoding/csv"
	"io"
	"log"
	"strconv"
)

func GenerateCSVEntry(w io.Writer) {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write([]string{"Id", "Record_id", "Data", "Created_at"})
	entries := db_.SelectEntry()
	rows := [][]string{}
	for _, entry := range entries {
		row := []string{}
		id := strconv.Itoa(int(entry.Id))
		recordID := strconv.Itoa(int(entry.Record_id))
		createdAT := entry.Created_at.Format("2006-01-02 15:04:05")
		data := string(entry.Data)
		row = append(row, id)
		row = append(row, recordID)
		row = append(row, createdAT)
		row = append(row, data)
		rows = append(rows, row)
	}
	if err := writer.WriteAll(rows); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
