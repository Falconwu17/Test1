package variables

import (
	"database/sql"
	"log"
)

func ConnectToDB() *sql.DB {
	db, err := sql.Open("postgres",
		"host=localhost user=postgres password=password dbname=students sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
