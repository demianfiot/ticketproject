package service

import (
	"context"

	aiclient "github.com/demianfiot/ticketproject/ticket-service/internal/client/ai"
	"github.com/demianfiot/ticketproject/ticket-service/internal/repository"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service/events"
)

type Ticket interface {
	CreateTicket(ctx context.Context, in CreateTicketInput) (CreateTicketResult, error)
	GetTicketByID(ctx context.Context, id int) (GetTicketResult, error)
}

type Service struct {
	Ticket
}

func NewService(
	aiClient *aiclient.GRPCClient,
	repo *repository.Repository,
	producer events.Producer,
) *Service {
	return &Service{
		Ticket: NewTicketService(aiClient, repo.Ticket, producer),
	}
}
