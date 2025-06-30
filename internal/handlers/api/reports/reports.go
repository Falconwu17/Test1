package reports

import (
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/reports"
	"net/http"
)

func InitRoutesReport() {
	base.RegisterRoute(base.NewRoute("GET", "/reports/entries", getEntryReports))
	base.RegisterRoute(base.NewRoute("GET", "/reports/records", getRecordReport))
}

func getEntryReports(w http.ResponseWriter, r *http.Request) {
	reports.CsvEntryHandler()(w, r)
}
func getRecordReport(w http.ResponseWriter, r *http.Request) {
	reports.CsvRecordsHandler()(w, r)
}
