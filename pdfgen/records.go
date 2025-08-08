package pdfgen

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"awesomeProject1/utils"
	"bytes"
	"net/http"
	"strconv"
)

func PdfHandlerForRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		limit, offset := utils.ParseLimitOffset(request)
		records, err := db_.SelectRecord(limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pdfBuffer, err := GenPdfForRecords(records)
		if err != nil {
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=records_report.pdf")
		w.Write(pdfBuffer.Bytes())
	}
}
func GenPdfForRecords(records []models.Record) (*bytes.Buffer, error) {
	pdf := NewPdf("Records Report")
	rows := make([][]string, 0, len(records))
	for _, record := range records {
		rows = append(rows, []string{
			strconv.Itoa(record.RecordId),
			strconv.Itoa(record.Timeout),
			record.CreatedAt.Format("2006-01-02 15:04:05"),
			record.Status,
		})
	}
	spec := models.TableSpec{
		Title:  "Records Report",
		Header: []string{"ID", "Timeout", "CreatedAt", "Status"},
		Width:  []float64{30, 30, 60, 30},
		Rows:   rows,
	}
	RenderTable(pdf, spec)
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, err
}
