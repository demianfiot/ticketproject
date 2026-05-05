package llm

import (
	"context"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
)

type Client interface {
	AnalyzeTicket(ctx context.Context, input domain.AnalyzeInput) (domain.AnalyzeResult, error)
}
