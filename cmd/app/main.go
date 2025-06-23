package main

import (
	"awesomeProject1/internal/db"
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"fmt"
	"net/http"
)

func main() {
	handlers_.InitHandlers()
	db_.InsertStudent(models.Student{Student_id: 1, Name: "Sagir", Age: 20, Curs: 4})

	fmt.Println("Слушаю тебя на 8080")
	http.ListenAndServe(":8080", nil)
}
