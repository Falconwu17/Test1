package reports

import (
	"awesomeProject1/reports/entriesCSV"
	"awesomeProject1/reports/recordsCSV"
	"awesomeProject1/utils"
	"net/http"
)

type CsvHandler interface {
	http.Handler
}

func withCSVHeaders(handler func(http.ResponseWriter, *http.Request, int, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		limit, offset := utils.ParseLimitOffset(r)
		handler(w, r, limit, offset)
	}
}

func CsvEntryHandler() http.HandlerFunc {
	return withCSVHeaders(func(w http.ResponseWriter, r *http.Request, limit, offset int) {
		entriesCSV.GenerateCSVEntry(w, limit, offset)
	})
}

func CsvRecordsHandler() http.HandlerFunc {
	return withCSVHeaders(func(w http.ResponseWriter, r *http.Request, limit, offset int) {
		recordsCSV.GenerateCSVRecords(w, limit, offset)
	})
}
