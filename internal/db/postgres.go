package db_

import (
	""
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	_ "github.com/lib/pq"
	"log"
)

func InsertStudent(student models.Student) {
	db := variables.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("Insert Into students (students_id, name,age,curs ) Values ($1, $2, $3, $4)",
		student.Student_id, student.Name, student.Age, student.Curs)
	if err != nil {
		log.Fatal(err)
		return
	}
}
func SelectStudent() []models.Student {
	students := []models.Student{}
	db := variables.ConnectToDB()
	defer db.Close()
	rows, err := db.Query(`Select students_id, name, age, curs from students`)
	if err != nil {
		log.Fatal(err)
		return students
	}
	defer rows.Close()
	for rows.Next() {
		var student models.Student
		var students_id, age, curs int
		var name string

		err = rows.Scan(&students_id, &name, &age, &curs)
		if err != nil {
			log.Fatal(err)
			return students
		}

		student.Student_id = students_id
		student.Name = name
		student.Age = age
		student.Curs = curs
		students = append(students, student)
	}
	return students
}
