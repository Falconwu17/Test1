package handlers_

import (
	"awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"encoding/json"
	"net/http"
)

func InitHandlers() {
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			students := db_.SelectStudent()
			json.NewEncoder(w).Encode(students)
		} else if r.Method == "POST" {
			student := models.Student{}
			err := json.NewDecoder(r.Body).Decode(&student)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			db_.InsertStudent(student)
			if err := json.NewEncoder(w).Encode(student); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Метод не найден", http.StatusMethodNotAllowed)
		}
	})

}
