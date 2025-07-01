package email

import "net/http"

func GetHandlerMail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		SendEmail()
	}
}
