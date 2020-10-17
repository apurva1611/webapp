package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func produce(kafkaURL, topic string, watch WATCH, msgKey string) {
	writer := newKafkaWriter(kafkaURL, topic)
	watchJSON, _ := json.Marshal(watch)
	msg := kafka.Message{
		Key:   []byte(msgKey),
		Value: watchJSON,
	}

	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	}
}
