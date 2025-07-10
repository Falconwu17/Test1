package entry

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

func InitRoutesEntry() {
	base.RegisterRoute(base.NewRoute("GET", "/entries", getAll))
	base.RegisterRoute(base.NewRoute("POST", "/entries", create))
	base.RegisterRoute(base.NewRoute("DELETE", "/entries", delete))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, LimitErr := strconv.Atoi(limitStr)
	offset, offsetErr := strconv.Atoi(offsetStr)
	if LimitErr != nil || limit < 0 || limit > 100 {
		limit = 100
	}
	if offsetErr != nil || offset < 0 {
		offset = 0
	}
	entries, _ := db_.SelectEntry(limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	entries := models.Entry{}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(entries.Data) == 0 || string(entries.Data) == "null" {
		defaultData, _ := json.Marshal(map[string]string{"key": "value"})
		entries.Data = defaultData
	}

	entries.SetDefaultEntry()
	validate := variables.Validator
	err := validate.Struct(entries)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}
	exists, err := db_.CheckRecordExists(entries.Record_id)
	if err != nil {
		http.Error(w, "DB error during record check", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Referenced record_id does not exist", http.StatusBadRequest)
		return
	}
	if err := db_.InsertEntry(entries); err != nil {
		http.Error(w, "Failed to insert", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	entriesID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db_.DeleteEntryById(entriesID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "entry deleted"})
}
