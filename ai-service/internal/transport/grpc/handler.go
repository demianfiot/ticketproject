package transport

import (
	"context"

	ai "github.com/demianfiot/ticketproject/ai-service/proto"
)

type Handler struct {
	ai.UnimplementedAIServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (s *Handler) AnalyzeTicket(ctx context.Context, req *ai.AnalyzeTicketRequest) (*ai.AnalyzeTicketResponse, error) {
	return &ai.AnalyzeTicketResponse{
		Summary:        "Stub summary for: " + req.Title,
		Category:       "billing",
		Priority:       "high",
		SuggestedReply: "We are sorry for the inconvenience. Our team is checking your issue.",
	}, nil
}
