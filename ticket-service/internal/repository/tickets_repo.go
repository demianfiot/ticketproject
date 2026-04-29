package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	domain "github.com/demianfiot/ticketproject/ticket-service/internal"
)

type TicketPostgres struct {
	db *sqlx.DB
}

func NewTicketPostgres(db *sqlx.DB) *TicketPostgres {
	return &TicketPostgres{db: db}
}

func (r *TicketPostgres) CreateTicket(ctx context.Context, ticket domain.Ticket) (int, error) {
	var id int

	query := `
		INSERT INTO tickets (user_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		ticket.UserID,
		ticket.Title,
		ticket.Description,
		ticket.Status,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create ticket: %w", err)
	}

	return id, nil
}

func (r *TicketPostgres) UpdateTicketAIAnalysis(ctx context.Context, ticket domain.Ticket) error {
	query := `
		UPDATE tickets
		SET ai_summary = $1,
			ai_category = $2,
			ai_priority = $3,
			ai_suggested_reply = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ticket.AISummary,
		ticket.AICategory,
		ticket.AIPriority,
		ticket.AISuggestedReply,
		ticket.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update ticket ai analysis: %w", err)
	}

	return nil
}

func (r *TicketPostgres) GetTicketByID(ctx context.Context, id int) (domain.Ticket, error) {
	query := `
		SELECT id, user_id, title, description, status,
		       ai_summary, ai_category, ai_priority, ai_suggested_reply,
		       created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	var ticket domain.Ticket

	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		return domain.Ticket{}, fmt.Errorf("failed to get ticket: %w", err)
	}

	return ticket, nil
}
