package records

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func InitRoutesRecords() {
	base.RegisterRoute(base.NewRoute("GET", "/records", getAll))
	base.RegisterRoute(base.NewRoute("POST", "/records", create))
	base.RegisterRoute(base.NewRoute("DELETE", "/records", delete))
}
func getAll(w http.ResponseWriter, r *http.Request) {
	records := db_.SelectRecord()
	json.NewEncoder(w).Encode(records)
}
func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	records := models.Record{}
	if err := json.NewDecoder(r.Body).Decode(&records); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db_.InsertRecord(records)
	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	recordID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db_.DeleteRecordById(recordID)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "record deleted"}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
