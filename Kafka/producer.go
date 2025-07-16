package Kafka

import (
	"awesomeProject1/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"strconv"
)

func SendEntryMessageKafka(entry models.Entry) {
	ctx := context.Background()

	broker := os.Getenv("KAFKA_BROKER")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    os.Getenv("KAFKA_TOPIC_ENTRIES"),
		Balancer: &kafka.Hash{},
	})
	defer writer.Close()

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}

	key := strconv.Itoa(entry.RecordId)

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: jsonEntry,
	})
	if err != nil {
		fmt.Println("Ошибка отправки сообщения в Kafka:", err)
		return
	}

	fmt.Println("Сообщение отправлено в Kafka:", key)
}
