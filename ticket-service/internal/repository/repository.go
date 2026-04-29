package repository

import (
	"context"

	domain "github.com/demianfiot/ticketproject/ticket-service/internal"
	"github.com/jmoiron/sqlx"
)

type Ticket interface {
	CreateTicket(ctx context.Context, ticket domain.Ticket) (int, error)
	UpdateTicketAIAnalysis(ctx context.Context, ticket domain.Ticket) error
	GetTicketByID(ctx context.Context, id int) (domain.Ticket, error)
}

type Repository struct {
	Ticket
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Ticket: NewTicketPostgres(db),
	}
}
