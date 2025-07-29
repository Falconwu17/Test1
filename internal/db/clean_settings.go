package db_

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
)

func GetSettingByUserID(userID int) (models.AutoCleanSetting, error) {
	var setting models.AutoCleanSetting
	db := variables.DB
	err := db.QueryRow("SELECT id, user_id, enabled, interval_seconds, last_cleaned_at FROM auto_clean_settings WHERE user_id=$1", userID).Scan(
		&setting.ID, &setting.UserID, &setting.Enabled, &setting.IntervalSeconds, &setting.LastCleanedAt,
	)
	return setting, err
}

func UpdateSetting(setting models.AutoCleanSetting) error {
	db := variables.DB
	_, err := db.Exec(`
		UPDATE auto_clean_settings
		SET enabled = $1,
		    interval_seconds = $2,
		    last_cleaned_at = $3
		WHERE user_id = $4`,
		setting.Enabled, setting.IntervalSeconds, setting.LastCleanedAt, setting.UserID,
	)
	if err != nil {
		log.Printf("Ошибка при обновлении настройки: %v", err)
	}
	log.Printf("Настройка обновлена для user_id: %d", setting.UserID)

	return err
}
func GetAllAutoCleanSettings() ([]models.AutoCleanSetting, error) {
	var settings []models.AutoCleanSetting
	db := variables.DB
	rows, err := db.Query("SELECT id, user_id, enabled, interval_seconds FROM auto_clean_settings")
	if err != nil {
		log.Printf("DB query error: %v", err)
		return settings, err
	}
	defer rows.Close()
	for rows.Next() {
		var setting models.AutoCleanSetting
		if err := rows.Scan(&setting.ID, &setting.UserID, &setting.Enabled, &setting.IntervalSeconds); err != nil {
			log.Printf("DB scan error: %v", err)
			return settings, err
		}
		settings = append(settings, setting)
	}
	if err := rows.Err(); err != nil {
		log.Printf("DB query error: %v", err)
		return settings, err
	}
	return settings, nil
}
func InsertSetting(setting *models.AutoCleanSetting) error {
	db := variables.DB
	err := db.QueryRow(
		`INSERT INTO auto_clean_settings (user_id, enabled, interval_seconds, last_cleaned_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		setting.UserID, setting.Enabled, setting.IntervalSeconds, setting.LastCleanedAt,
	).Scan(&setting.ID)

	if err != nil {
		log.Printf("Ошибка при вставке auto_clean_setting: %v", err)
		return err
	}

	log.Printf("Настройка очистки добавлена с ID: %d", setting.ID)
	return nil
}
