package students

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"encoding/json"
	"net/http"
)

var routePrefix = "students"

func InitRoutesStudents() {
	base.RegisterRoute(base.NewRoute("GET", "/students", getAll))
	base.RegisterRoute(base.NewRoute("POST", "/students", create))
}
func getAll(w http.ResponseWriter, r *http.Request) {
	students := db_.SelectStudent()
	json.NewEncoder(w).Encode(students)
}
func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db_.InsertStudent(student)

	if err := json.NewEncoder(w).Encode(student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
