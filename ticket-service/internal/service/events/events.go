package events

import (
	"context"
	"time"
)

type TicketCreatedEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Version   int       `json:"version"`
	TicketID  int       `json:"ticket_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Producer interface {
	PublishTicketCreated(ctx context.Context, event TicketCreatedEvent) error
	Close() error
}
