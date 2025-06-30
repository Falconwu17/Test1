package reports

import (
	"awesomeProject1/reports/entriesCSV"
	"awesomeProject1/reports/recordsCSV"
	"net/http"
)

func CsvEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		entriesCSV.GenerateCSVEntry(w)
	}
}

func CsvRecordsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		recordsCSV.GenerateCSVRecords(w)
	}
}
