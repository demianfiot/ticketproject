package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	domain "github.com/demianfiot/ticketproject/ticket-service/internal"
	aiclient "github.com/demianfiot/ticketproject/ticket-service/internal/client/ai"
	"github.com/demianfiot/ticketproject/ticket-service/internal/repository"
	"github.com/demianfiot/ticketproject/ticket-service/internal/service/events"
)

type TicketService struct {
	repo     repository.Ticket
	aiClient *aiclient.GRPCClient
	producer events.Producer
}

type CreateTicketInput struct {
	UserID      string
	Title       string
	Description string
}

type CreateTicketResult struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	Summary        string `json:"summary"`
	Category       string `json:"category"`
	Priority       string `json:"priority"`
	SuggestedReply string `json:"suggested_reply"`
}

type GetTicketResult struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	Summary        string `json:"summary"`
	Category       string `json:"category"`
	Priority       string `json:"priority"`
	SuggestedReply string `json:"suggested_reply"`
}

func NewTicketService(
	aiClient *aiclient.GRPCClient,
	repo repository.Ticket,
	producer events.Producer,
) *TicketService {
	return &TicketService{
		aiClient: aiClient,
		repo:     repo,
		producer: producer,
	}
}

func (s *TicketService) CreateTicket(ctx context.Context, in CreateTicketInput) (CreateTicketResult, error) {
	ticket := domain.Ticket{
		UserID:      in.UserID,
		Title:       in.Title,
		Description: in.Description,
		Status:      "new",
	}

	id, err := s.repo.CreateTicket(ctx, ticket)
	if err != nil {
		return CreateTicketResult{}, fmt.Errorf("failed to create ticket in repository: %w", err)
	}

	ticket.ID = id
	analysis, err := s.aiClient.AnalyzeTicket(ctx, aiclient.AnalyzeInput{
		TicketID:    fmt.Sprintf("%d", ticket.ID),
		Title:       ticket.Title,
		Description: ticket.Description,
		UserID:      ticket.UserID,
	})

	if err != nil {
		log.Printf("failed to analyze ticket %d: %v", ticket.ID, err)

		ticket.AISummary = "AI analysis unavailable"
		ticket.AICategory = "unknown"
		ticket.AIPriority = "low"
		ticket.AISuggestedReply = "Support agent should review this ticket manually."
	} else {
		ticket.AISummary = analysis.Summary
		ticket.AICategory = analysis.Category
		ticket.AIPriority = analysis.Priority
		ticket.AISuggestedReply = analysis.SuggestedReply
	}

	if updateErr := s.repo.UpdateTicketAIAnalysis(ctx, ticket); updateErr != nil {
		log.Printf("failed to update ai analysis for ticket %d: %v", ticket.ID, updateErr)
	}

	event := events.TicketCreatedEvent{
		EventID:   uuid.NewString(),
		EventType: "ticket.created",
		Version:   1,
		TicketID:  ticket.ID,
		UserID:    ticket.UserID,
		Status:    ticket.Status,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.producer.PublishTicketCreated(ctx, event); err != nil {
		log.Printf("failed to publish ticket.created event for ticket %d: %v", ticket.ID, err)
	}

	return CreateTicketResult{
		ID:             ticket.ID,
		Title:          ticket.Title,
		Description:    ticket.Description,
		Status:         ticket.Status,
		Summary:        ticket.AISummary,
		Category:       ticket.AICategory,
		Priority:       ticket.AIPriority,
		SuggestedReply: ticket.AISuggestedReply,
	}, nil
}

func (s *TicketService) GetTicketByID(ctx context.Context, id int) (GetTicketResult, error) {
	ticket, err := s.repo.GetTicketByID(ctx, id)
	if err != nil {
		return GetTicketResult{}, fmt.Errorf("failed to get ticket: %w", err)
	}

	return GetTicketResult{
		ID:             ticket.ID,
		Title:          ticket.Title,
		Description:    ticket.Description,
		Status:         ticket.Status,
		Summary:        ticket.AISummary,
		Category:       ticket.AICategory,
		Priority:       ticket.AIPriority,
		SuggestedReply: ticket.AISuggestedReply,
	}, nil
}
