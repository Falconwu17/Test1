package pdfgen

import (
	"github.com/jung-kurt/gofpdf"
)

func NewPdf(title string) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(title, false)
	pdf.SetFont("Arial", "B", 12)
	pdf.AddPage()
	return pdf
}
