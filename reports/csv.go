package reports

import (
	"awesomeProject1/reports/recordsCSV"
	"log"
	"net/http"
	"strconv"
)

func CsvEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("<UNK> <UNK> <UNK>: %v", err)
		}
		offset, _ := strconv.Atoi(offsetStr)
		if limit <= 100 {
			limit = 100
		}
		if offset < 0 {
			offset = 0
		}
	}
}

func CsvRecordsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		limit, err := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)
		if limit <= 100 {
			limit = 100
		}
		if offset < 0 {
			offset = 0
		}
		if err != nil {
			log.Printf("Error for limit in CSV: %v", err)
		}

		recordsCSV.GenerateCSVRecords(w, limit, offset)
	}
}
