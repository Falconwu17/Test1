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
	"time"
)

func SendMessageInBase() {

	db := variables.DB
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpServer := models.SmtpServer{os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")}
	var buf bytes.Buffer
	recordsCSV.GenerateCSVRecords(&buf, 100, 0)
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
	subject := "Ежедневный отчёт по записям"
	body := buf.String()
	message := []byte("Subject: " + subject + "\r\n\r\n" + body)
	log.Printf("Отправляю с %s на %v", from, to)
	log.Printf("Тело письма:\n%s", body)
	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	log.Printf("Отправка на: %v\n", to)
	log.Printf("Тело письма: \n%s\n", buf.String())
	err = smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		log.Println("Ошибка отправки письма:", err)
		return
	}
	fmt.Println("Email Sent!")
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
