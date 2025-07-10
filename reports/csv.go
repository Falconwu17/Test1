package reports

import (
	"awesomeProject1/reports/entriesCSV"
	"awesomeProject1/reports/recordsCSV"
	"log"
	"net/http"
	"strconv"
)

type CsvHandler interface {
	http.Handler
}

func parseLimitOffset(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Printf("Error for limit in CSV: %v", err)
	}
	if err != nil || limit <= 0 || limit > 100 {
		limit = 100
	}
	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}
func withCSVHeaders(handler func(http.ResponseWriter, *http.Request, int, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv")
		limit, offset := parseLimitOffset(r)
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
