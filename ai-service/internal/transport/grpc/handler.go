package grpc

import (
	"context"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
	"github.com/demianfiot/ticketproject/ai-service/internal/service"
	aipb "github.com/demianfiot/ticketproject/ai-service/proto"
)

type Handler struct {
	aipb.UnimplementedAIServiceServer
	service *service.AnalysisService
}

func NewHandler(service *service.AnalysisService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) AnalyzeTicket(ctx context.Context, req *aipb.AnalyzeTicketRequest) (*aipb.AnalyzeTicketResponse, error) {
	result, err := h.service.AnalyzeTicket(ctx, domain.AnalyzeInput{
		TicketID:    req.TicketId,
		UserID:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &aipb.AnalyzeTicketResponse{
		Summary:        result.Summary,
		Category:       result.Category,
		Priority:       result.Priority,
		SuggestedReply: result.SuggestedReply,
	}, nil
}
