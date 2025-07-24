package main

import (
	"awesomeProject1/Kafka"
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/handlers/api/entry"
	"awesomeProject1/internal/handlers/api/records"
	"awesomeProject1/internal/handlers/api/reports"
	"awesomeProject1/internal/handlers/background"
	"awesomeProject1/variables"
	"fmt"
	"net/http"
)

func main() {
	go background.AutoSendMessage()
	go background.AutoCleanForTime()
	go Kafka.Consumer()
	variables.InitDB()
	variables.InitValidate()
	handlers_.InitBaseRoutes()
	records.InitRoutesRecords()
	entry.InitRoutesEntry()
	reports.InitRoutesReport()
	db_.CheckRecord()

	http.Handle("/", http.HandlerFunc(handlers_.ServeHTTP))
	fmt.Println("Слушаю тебя на 8080")
	http.ListenAndServe(":8080", nil)
}
