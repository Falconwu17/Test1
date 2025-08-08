package email

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"awesomeProject1/pdfgen"
	"awesomeProject1/reports/entriesCSV"
	"awesomeProject1/reports/recordsCSV"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func GetHandlerMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		postMail(w, r)
	}
}
func postMail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	request := models.EmailRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(request.To) == 0 {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := models.SmtpServer{os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")}
	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	var csvEntry bytes.Buffer
	var csvRecords bytes.Buffer
	records, _ := db_.SelectRecord(100, 0)
	entry, _ := db_.SelectEntry(100, 0)
	pdfRecords, _ := pdfgen.GenPdfForRecords(records)
	pdfEntry, _ := pdfgen.GenPdfForEntries(entry)
	recordsCSV.GenerateCSVRecords(&csvRecords, 100, 0)
	entriesCSV.GenerateCSVEntry(&csvEntry, 100, 0)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(request.To, ",")
	headers["Subject"] = request.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	htmlPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/html; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qp := quotedprintable.NewWriter(htmlPart)
	qp.Write([]byte(`<h2>Ежедневный отчёт</h2>
		<p>Во вложении CSV и PDF файлы с данными по метрикам.</p>`))
	qp.Close()

	rPart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":        {"text/csv"},
		"Content-Disposition": {"attachment; filename=\"records.csv\""},
	})
	rPart.Write(csvRecords.Bytes())

	ePart, _ := writer.CreatePart(map[string][]string{
		"Content-Type":        {"text/csv"},
		"Content-Disposition": {"attachment; filename=\"entries.csv\""},
	})
	ePart.Write(csvEntry.Bytes())
	pdfPart1, _ := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/pdf"},
		"Content-Disposition": {"attachment; filename=\"records.pdf\""},
	})
	pdfPart1.Write(pdfRecords.Bytes())
	pdfPart2, _ := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/pdf"},
		"Content-Disposition": {"attachment; filename=\"entries.pdf\""},
	})
	pdfPart2.Write(pdfEntry.Bytes())
	writer.Close()

	if err := smtp.SendMail(smtpServer.Address(), auth, from, request.To, append([]byte(msg.String()), body.Bytes()...)); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}
