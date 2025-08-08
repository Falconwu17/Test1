package pdfgen

import (
	"awesomeProject1/internal/models"
	"github.com/jung-kurt/gofpdf"
)

func RenderTable(pdf *gofpdf.Fpdf, spec models.TableSpec) {
	left := (210 - sum(spec.Width)) / 2
	pdf.SetFillColor(220, 220, 220)
	pdf.SetX(left)
	for i, h := range spec.Header {
		pdf.CellFormat(spec.Width[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	for _, row := range spec.Rows {
		pdf.SetX(left)
		for i, col := range row {
			pdf.CellFormat(spec.Width[i], 7, col, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
}
func sum(floats []float64) float64 {
	var total float64
	for _, v := range floats {
		total += v
	}
	return total
}
