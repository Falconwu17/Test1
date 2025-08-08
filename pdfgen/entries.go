package pdfgen

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"awesomeProject1/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func PdfHandlerForEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := utils.ParseLimitOffset(r)
		entries, err := db_.SelectEntry(limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pdfBuf, err := GenPdfForEntries(entries)
		if err != nil {
			http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=records_report.pdf")
		w.Write(pdfBuf.Bytes())
	}
}
func formatJson(data json.RawMessage) string {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return string(data) // если невалидный — просто вернём как есть
	}
	return out.String()
}

func GenPdfForEntries(entries []models.Entry) (*bytes.Buffer, error) {
	pdf := NewPdf("Entries Report")
	rows := make([][]string, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, []string{
			strconv.Itoa(entry.Id),
			strconv.Itoa(entry.RecordId),
			entry.CreatedAt.Format("2006-01-02 15:04:05"),
			formatJson(entry.Data),
		})

	}
	spec := models.TableSpec{
		Title:  "Entries Report",
		Header: []string{"ID", "Record ID", "Created At", "Data"},
		Width:  []float64{20, 30, 50, 110},
		Rows:   rows,
	}
	RenderTable(pdf, spec)
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		fmt.Println(err)
	}
	return &buf, nil
}
