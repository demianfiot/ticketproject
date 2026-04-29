package repository

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/demianfiot/ticketproject/ticket-service/internal/service/events"
)

type AnalyticsRepository struct {
	conn driver.Conn
}

func NewAnalyticsRepository(conn driver.Conn) *AnalyticsRepository {
	return &AnalyticsRepository{
		conn: conn,
	}
}

func (r *AnalyticsRepository) InsertTicketCreatedEvent(ctx context.Context, event events.TicketCreatedEvent) error {
	query := `
		INSERT INTO ticket_events
		(event_id, event_type, version, ticket_id, user_id, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	if err := r.conn.Exec(
		ctx,
		query,
		event.EventID,
		event.EventType,
		event.Version,
		event.TicketID,
		event.UserID,
		event.Status,
		event.CreatedAt,
	); err != nil {
		return fmt.Errorf("failed to insert ticket event: %w", err)
	}

	return nil
}
