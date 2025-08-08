package auto

import (
	db_ "awesomeProject1/internal/db"
	base "awesomeProject1/internal/handlers"
	"awesomeProject1/internal/middleware"
	"awesomeProject1/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func InitAutoSetting() {
	base.RegisterProtectedRoute(base.NewRoute("GET", "/autoGet", getALL))
	base.RegisterProtectedRoute(base.NewRoute("GET", "/autoGetByUser", getByUSerID))
	base.RegisterProtectedRoute(base.NewRoute("POST", "/autoPOST", create))
	base.RegisterProtectedRoute(base.NewRoute("PUT", "/autoUpdate", updateSetting))
}
func getALL(w http.ResponseWriter, r *http.Request) {
	settings, err := db_.GetAllAutoCleanSettings()
	if err != nil {
		log.Printf("DB GetAllAutoCleanSettings error : %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		log.Printf("DB GetAllAutoCleanSettings error Encode: %v", err)
		return
	}
}
func getByUSerID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	setting, err := db_.GetSettingByUserID(userID)
	if err != nil {
		http.Error(w, "Настройка не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(setting)
}
func updateSetting(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	setting := models.AutoCleanSetting{}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}
	setting.UserID = userID
	exists, err := db_.CheckUserExists(userID)
	if err != nil {
		http.Error(w, "DB error during record check", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Referenced record_id does not exist", http.StatusBadRequest)
		return
	}
	if er := db_.UpdateSetting(setting); er != nil {
		if er.Error() == fmt.Sprintf("no setting found for user_id: %d", setting.UserID) {
			http.Error(w, er.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при обновлении настройки", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		http.Error(w, "Sorry JSON", http.StatusInternalServerError)
		return
	}

}
func create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	setting := models.AutoCleanSetting{}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}
	setting.UserID = userID
	if err := db_.InsertSetting(&setting); err != nil {
		http.Error(w, "Ошибка при вставке настройки", http.StatusBadRequest)
		log.Printf("Ошибка вставки: %v", err)
		return
	}
	log.Printf("Настройка успешно добавлена: %+v", setting)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(setting); err != nil {
		http.Error(w, "Ошибка при ответе", http.StatusInternalServerError)
	}

}
