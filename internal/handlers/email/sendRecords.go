package email

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/reports/recordsCSV"
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

func SmtpServerAddress(s models.SmtpServer) string {
	return s.Host + ":" + s.Port
}

func SendEmail() {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{
		"mukan@gmail.com",
		"alex@yandex.ru",
	}
	smtpServer := models.SmtpServer{os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT")}
	var buf bytes.Buffer
	recordsCSV.GenerateCSVRecords(&buf, 100, 0)
	message := []byte("Subject: Report\r\n\r\n" + buf.String())

	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	if err := smtp.SendMail(SmtpServerAddress(smtpServer), auth, from, to, message); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Email sent successfully")
}
