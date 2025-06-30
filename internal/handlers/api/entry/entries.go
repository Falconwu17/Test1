package entry

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func InitRoutesEntry() {
	base.RegisterRoute(base.NewRoute("GET", "/entries", getAll))
	base.RegisterRoute(base.NewRoute("POST", "/entries", create))
	base.RegisterRoute(base.NewRoute("DELETE", "/entries", delete))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	entry := db_.SelectEntry()
	json.NewEncoder(w).Encode(entry)
}
func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	entries := models.Entry{}

	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db_.InsertEntry(entries)
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
	}
	db_.DeleteEntryById(entriesID)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "entry deleted"}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
