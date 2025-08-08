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
	"strings"
)

func InitRoutesRecords() {
	base.RegisterProtectedRoute(base.NewRoute("GET", "/records", getAll))
	base.RegisterProtectedRoute(base.NewRoute("POST", "/records", create))
	base.RegisterProtectedRoute(base.NewRoute("DELETE", "/records", delete))
	base.RegisterProtectedRoute(base.NewRoute("PUT", "/records", update))
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
	record.Status = strings.ToLower(record.Status)

	validStatuses := map[string]bool{
		"now":     true,
		"later":   true,
		"process": true,
	}
	if !validStatuses[record.Status] {
		http.Error(w, "Invalid status value. Allowed: now, later, process", http.StatusBadRequest)
		return
	}
	validate := variables.Validator
	err := validate.Struct(record)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}
	if err := db_.InsertRecord(&record); err != nil {
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
func update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	var record models.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if record.RecordId == 0 {
		http.Error(w, "Missing or invalid record ID", http.StatusBadRequest)
		return
	}

	record.Status = strings.ToLower(record.Status)

	validStatuses := map[string]bool{
		"now":     true,
		"later":   true,
		"process": true,
	}
	if !validStatuses[record.Status] {
		http.Error(w, "Invalid status value. Allowed: now, later, process", http.StatusBadRequest)
		return
	}

	validate := variables.Validator
	if err := validate.Struct(record); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}
	exists, err := db_.CheckRecordExists(record.RecordId)
	if err != nil {
		http.Error(w, "DB error during record check", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	if err := db_.UpdateRecord(&record); err != nil {
		log.Printf("[ERROR] failed to update record: %v", err)
		http.Error(w, "Failed to update record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "record updated"})
}
