package models

import (
	"encoding/json"
	"time"
)

type TableSpec struct {
	Title  string     `json:"title"`
	Header []string   `json:"header"`
	Width  []float64  `json:"width"`
	Rows   [][]string `json:"rows"`
}
type Record struct {
	RecordId  int       `json:"record_id"`
	Timeout   int       `json:"timeout" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" `
	Status    string    `json:"status" validate:"required,oneof=now later process"`
}
type Entry struct {
	Id        int             `json:"id"`
	RecordId  int             `json:"record_id" validate:"required,min=1"`
	Data      json.RawMessage `json:"data" validate:"required"`
	CreatedAt time.Time       `json:"created_at" validate:"required"`
}

type SmtpServer struct {
	Host string
	Port string
}

func (s *SmtpServer) Address() string {
	return s.Host + ":" + s.Port
}

type EmailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AutoCleanSetting struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Enabled         bool      `json:"enabled"`
	IntervalSeconds int       `json:"interval_seconds"`
	LastCleanedAt   time.Time `json:"last_cleaned_at"`
}

func (r *Record) SetDefault() {
	if r.Timeout == 0 {
		r.Timeout = 60
	}
	if r.Status == "" {
		r.Status = "Now"
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
}
func (r *Entry) SetDefaultEntry() {
	if r.RecordId == 0 {
		r.RecordId = 1
	}
	if len(r.Data) == 0 {
		defaultData, _ := json.Marshal(map[string]string{"key": "value"})

		r.Data = defaultData
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
}
