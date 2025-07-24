package background

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/reports/recordsCSV"
	"awesomeProject1/variables"
	"bytes"
	"database/sql"
	"fmt"
	"log"
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
	var buf bytes.Buffer
	recordsCSV.GenerateCSVRecords(&buf, 100, 0)
	rows, err := db.Query("SELECT email FROM user")
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
	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: 🧾 Ежедневный отчёт по записям\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		from,
		strings.Join(to, ", "),
		buf.String(),
	))

	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		log.Println("Ошибка отправки письма:", err)
		return
	}
	fmt.Println("Email Sent!")
}
func AutoSendMessage() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				SendMessageInBase()
			}
		}
	}()
}
