package events

import (
	"context"
	"encoding/json"
	"strconv"

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

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(event.TicketID)),
		Value: payload,
		Headers: []kafka.Header{
			{
				Key:   "event_type",
				Value: []byte(event.EventType),
			},
		},
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
