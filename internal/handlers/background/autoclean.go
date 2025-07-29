package background

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"awesomeProject1/variables"
	"log"
	"time"
)

func AutoClean(settings models.AutoCleanSetting) {
	db := variables.DB
	if settings.Enabled != true {
		log.Printf("Автоочистка отключена для пользователя %d", settings.UserID)
		return
	}
	_, err := db.Exec(`
		DELETE FROM records 
		WHERE user_id = $1 
		  AND created_at < NOW() - INTERVAL '1 second' * $2`,
		settings.UserID, settings.IntervalSeconds,
	)

	if err != nil {
		log.Printf("Ошибка при удалении записей для user_id %d: %v", settings.UserID, err)
		return
	}
	_, err = db.Exec(`
    UPDATE auto_clean_settings 
    SET last_cleaned_at = NOW() 
    WHERE user_id = $1`, settings.UserID)
	if err != nil {
		log.Printf("Ошибка при обновление информации ", settings.UserID, err)
	}
	log.Printf("Автоочистка завершена для user_id %d", settings.UserID)

}

func AutoCleanForTime() {
	ticker := time.NewTicker(1 * time.Hour)
	done := make(chan bool)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				settings, err := db_.GetAllAutoCleanSettings()
				if err != nil {
					log.Printf("Ошибка получения настроек автоочистки: %v", err)
					continue
				}
				for _, setting := range settings {
					AutoClean(setting)
				}
				log.Printf("Запущена автоочистка для %d пользователей", len(settings))
			}
		}
	}()

}
