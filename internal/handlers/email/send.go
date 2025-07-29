package email

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/reports/recordsCSV"
	"bytes"
	"encoding/json"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
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
	smtpServer.Address()
	var buf bytes.Buffer
	recordsCSV.GenerateCSVRecords(&buf, 100, 0)
	message := []byte("Subject: " + request.Subject + "\r\n\r\n" + request.Body + "\n\n" + buf.String())

	for _, addr := range request.To {
		if _, err := mail.ParseAddress(addr); err != nil {
			http.Error(w, "Invalid email: "+addr, http.StatusBadRequest)
			return
		}
	}
	auth := smtp.PlainAuth("", from, password, smtpServer.Host)
	if err := smtp.SendMail(smtpServer.Address(), auth, from, request.To, message); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
	})

}
