package background

import (
	"awesomeProject1/variables"
	"log"
	"time"
)

func AutoClean() {
	db := variables.DB
	_, err := db.Exec("DELETE FROM records WHERE created_at < NOW() - INTERVAL '7 days'")
	if err != nil {
		log.Printf("Ошибка при удалении: %v", err)
	} else {
		log.Println("Очистка завершена: удалены записи старше 7 дней")
	}
}

func AutoCleanForTime() {
	ticker := time.NewTicker(168 * time.Hour)
	done := make(chan bool)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				AutoClean()
			}
		}
	}()
}
