package models

import (
	"encoding/json"
	"time"
)

type Student struct {
	Student_id int    `json:"student_id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Curs       int    `json:"curs"`
}

type Record struct {
	Record_id    int       `json:"record_id"`
	Timeout      int       `json:"timeout"`
	Handler_type string    `json:"handler_type"`
	Created_at   time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

type Entry struct {
	Id         int             `json:"id"`
	Record_id  int             `json:"record_id"`
	Data       json.RawMessage `json:"data"`
	Created_at time.Time       `json:"created_at"`
}
type SmtpServer struct {
	Host string
	Port string
}
