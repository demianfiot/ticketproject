package events

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) PublishTicketCreated(ctx context.Context, event TicketCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(strconv.Itoa(event.TicketID)),
		Value: payload,
		Headers: []kafka.Header{
			{
				Key:   "event_type",
				Value: []byte(event.EventType),
			},
		},
	}

	var lastErr error

	for attempt := 1; attempt <= 3; attempt++ {
		writeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)

		err = p.writer.WriteMessages(writeCtx, msg)

		cancel()

		if err == nil {
			return nil
		}

		lastErr = err
		time.Sleep(300 * time.Millisecond)
	}

	return lastErr
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
