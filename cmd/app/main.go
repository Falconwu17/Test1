package main

import (
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/handlers/api/entry"
	"awesomeProject1/internal/handlers/api/records"
	"awesomeProject1/internal/handlers/api/reports"
	"awesomeProject1/variables"
	"fmt"
	"net/http"
)

func main() {
	handlers_.InitBaseRoutes()
	records.InitRoutesRecords()
	entry.InitRoutesEntry()
	reports.InitRoutesReport()
	variables.InitDB()

	http.Handle("/", http.HandlerFunc(handlers_.ServeHTTP))
	fmt.Println("Слушаю тебя на 8080")
	http.ListenAndServe(":8080", nil)
}
