package background

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/reports/entriesCSV"
	"awesomeProject1/reports/recordsCSV"
	"awesomeProject1/variables"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func SendMessageInBase() {
	db := variables.DB
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := models.SmtpServer{os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")}
	auth := smtp.PlainAuth("", from, password, smtpServer.Host)

	rows, err := db.Query("SELECT email FROM users")
	if err != nil {
		log.Println("Ошибка при получении email:", err)
		return
	}
	defer rows.Close()

	var to []string
	for rows.Next() {
		var email sql.NullString
		if err := rows.Scan(&email); err == nil && email.Valid {
			to = append(to, email.String)
		}
	}
	if len(to) == 0 {
		log.Println("Нет email для отправки")
		return
	}

	var csvRecords bytes.Buffer
	var csvEntries bytes.Buffer
	recordsCSV.GenerateCSVRecords(&csvRecords, 100, 0)
	entriesCSV.GenerateCSVEntry(&csvEntries, 100, 0)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(to, ","),
		"Subject":      "Ежедневный отчёт по записям",
		"MIME-Version": "1.0",
		"Content-Type": "multipart/mixed; boundary=" + writer.Boundary(),
	}

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
	qp.Write([]byte(`<h3>Автоматический отчёт</h3><p>Во вложении два CSV-файла: <b>records.csv</b> и <b>entries.csv</b>.</p>`))
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
	ePart.Write(csvEntries.Bytes())

	writer.Close()

	fullMsg := append([]byte(msg.String()), body.Bytes()...)
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, fullMsg)
	if err != nil {
		log.Println("Ошибка отправки письма:", err)
		return
	}
	log.Printf("Email успешно отправлен на %v", to)
}

func AutoSendMessage() {
	log.Println("Автоматическая отправка стартовала")
	ticker := time.NewTicker(24 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Попытка отправки письма...")
				SendMessageInBase()
			}
		}
	}()
}
