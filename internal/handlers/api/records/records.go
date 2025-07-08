package records

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

func InitRoutesRecords() {
	base.RegisterRoute(base.NewRoute("GET", "/records", getAll))
	base.RegisterRoute(base.NewRoute("POST", "/records", create))
	base.RegisterRoute(base.NewRoute("DELETE", "/records", delete))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err1 := strconv.Atoi(limitStr)
	offset, err2 := strconv.Atoi(offsetStr)
	if err1 != nil || limit <= 0 || limit > 100 {
		limit = 100
	}
	if err2 != nil || offset < 0 {
		offset = 0
	}

	records, err := db_.SelectRecord(limit, offset)
	if err != nil {
		log.Printf("[ERROR] failed to select records: %v", err)
		http.Error(w, "Failed to get records", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func create(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	record := models.Record{}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	record.SetDefault()
	validate := variables.Validator
	err := validate.Struct(record)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	if err := db_.InsertRecord(record); err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(record)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	recordID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := db_.DeleteRecordById(recordID); err != nil {
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "record deleted"})
}
