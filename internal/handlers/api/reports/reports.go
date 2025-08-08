package reports

import (
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/handlers/email"
	"awesomeProject1/pdfgen"
	"awesomeProject1/reports"
)

func InitRoutesReport() {
	base.RegisterProtectedRoute(base.NewRoute("GET", "/reports/entries", reports.CsvEntryHandler()))
	base.RegisterProtectedRoute(base.NewRoute("GET", "/reports/records", reports.CsvRecordsHandler()))
	base.RegisterRoute(base.NewRoute("GET", "/reports/graph", reports.GraphHandler()))
	base.RegisterProtectedRoute(base.NewRoute("POST", "/reports/send-mail", email.GetHandlerMail()))
	base.RegisterProtectedRoute(base.NewRoute("GET", "/reports/pdf/record", pdfgen.PdfHandlerForRecord()))
	base.RegisterProtectedRoute(base.NewRoute("GET", "/reports/pdf/entry", pdfgen.PdfHandlerForEntry()))
}
