package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func kafkaHealthCheck(kafkaURL string) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}

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

	log.Printf("PRODUCE Topic: %s, Message ID %s", topic, string(msg.Key))
}

func producetest(kafkaURL, topic string, watch string, msgKey string) {
	writer := newKafkaWriter(kafkaURL, topic)
	msg := kafka.Message{
		Key:   []byte(msgKey),
		Value: []byte(watch),
	}

	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	}
}
