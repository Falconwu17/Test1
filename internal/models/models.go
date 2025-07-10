package models

import (
	"encoding/json"
	"time"
)

type Record struct {
	Record_id  int       `json:"record_id"`
	Timeout    int       `json:"timeout" validate:"required,min=1"`
	Created_at time.Time `json:"created_at" validate:"required"`
	Status     string    `json:"status" validate:"required,oneof=Now Later Process"`
}
type Entry struct {
	Id         int             `json:"id"`
	Record_id  int             `json:"record_id" validate:"required,min=1"`
	Data       json.RawMessage `json:"data" validate:"required"`
	Created_at time.Time       `json:"created_at" validate:"required"`
}

type SmtpServer struct {
	Host string
	Port string
}
type SetDefault interface {
	SetDefault()
	SetDefaultEntry()
}

func (r *Record) SetDefault() {
	if r.Record_id == 0 {
		r.Record_id = 1
	}
	if r.Timeout == 0 {
		r.Timeout = 60
	}
	if r.Status == "" {
		r.Status = "Now"
	}
	if r.Created_at.IsZero() {
		r.Created_at = time.Now()
	}
}
func (r *Entry) SetDefaultEntry() {
	if r.Id == 0 {
		r.Id = 1
	}
	if r.Record_id == 0 {
		r.Record_id = 1
	}
	if len(r.Data) == 0 {
		defaultData, _ := json.Marshal(map[string]string{"key": "value"})

		r.Data = defaultData
	}
	if r.Created_at.IsZero() {
		r.Created_at = time.Now()
	}
}
