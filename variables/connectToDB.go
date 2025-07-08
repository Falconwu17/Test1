package variables

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

var Validator *validator.Validate

func InitValidate() {
	Validator = validator.New()
	return
}
