package reports

import (
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/handlers/email"
	"awesomeProject1/reports"
	"net/http"
)

func InitRoutesReport() {
	base.RegisterRoute(base.NewRoute("GET", "/reports/entries", getEntryReports))
	base.RegisterRoute(base.NewRoute("GET", "/reports/records", getRecordReport))
	base.RegisterRoute(base.NewRoute("GET", "/reports/records/graph", getRecordGraph))
	base.RegisterRoute(base.NewRoute("GET", "/reports/send-mail", getRecordsMail))
}
func getRecordsMail(w http.ResponseWriter, r *http.Request) {
	email.GetHandlerMail().ServeHTTP(w, r)
}
func getRecordGraph(w http.ResponseWriter, r *http.Request) {
	reports.GraphHandler()(w, r)
}
func getEntryReports(w http.ResponseWriter, r *http.Request) {
	reports.CsvEntryHandler()(w, r)
}
func getRecordReport(w http.ResponseWriter, r *http.Request) {
	reports.CsvRecordsHandler()(w, r)
}
