package handlers_

import (
	"awesomeProject1/internal/db"
	"encoding/json"
	"net/http"
)

func InitHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			students := db_.SelectStudent()
			json.NewEncoder(w).Encode(students)
		} else if r.Method == "POST" {
			//В будущем что то да и сделаю

		}
	})
}
