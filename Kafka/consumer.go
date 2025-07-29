package Kafka

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
)

func Consumer() {
	ctx := context.Background()
	brokers := os.Getenv("KAFKA_BROKER")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		GroupID: "entries-consumers",
		Topic:   os.Getenv("KAFKA_TOPIC_ENTRIES"),
	})
	defer r.Close()
	for {
		fmt.Println("Kafka Consumer стартовал")
		fmt.Printf("Подключение к брокеру: %s, топик: %s\n", brokers, os.Getenv("KAFKA_TOPIC_ENTRIES"))
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Println("Ошибка чтения сообщения из Kafka:", err)
			continue
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		var entry models.Entry
		if err := json.Unmarshal(m.Value, &entry); err != nil {
			fmt.Println("Ошибка декодирования JSON:", err)
			continue
		}
		if err := db_.InsertEntry(&entry); err != nil {
			fmt.Println("Ошибка вставки в бд:", err)
		}
	}

}
