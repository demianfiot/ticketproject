package service

import (
	"context"
	"fmt"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
	"github.com/demianfiot/ticketproject/ai-service/internal/llm"
)

type AnalysisService struct {
	llmClient llm.Client
}

func NewAnalysisService(llmClient llm.Client) *AnalysisService {
	return &AnalysisService{
		llmClient: llmClient,
	}
}

func (s *AnalysisService) AnalyzeTicket(ctx context.Context, input domain.AnalyzeInput) (domain.AnalyzeResult, error) {
	if input.Title == "" || input.Description == "" {
		return domain.AnalyzeResult{}, fmt.Errorf("title and description are required")
	}

	result, err := s.llmClient.AnalyzeTicket(ctx, input)
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("llm analyze ticket failed: %w", err)
	}

	return result, nil
}
