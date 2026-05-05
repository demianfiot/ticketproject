package prompt

import (
	"fmt"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildTicketAnalysisPrompt(input domain.AnalyzeInput) string {
	return fmt.Sprintf(`
You are an AI assistant for a customer support platform.

Analyze the following support ticket and return ONLY valid JSON.

Allowed categories:
- billing
- technical
- account
- refund
- other

Allowed priorities:
- low
- medium
- high
- urgent

Required JSON format:
{
  "summary": "short summary of the issue",
  "category": "billing|technical|account|refund|other",
  "priority": "low|medium|high|urgent",
  "suggested_reply": "short professional reply to the customer",
  "confidence": 0.0
}

Ticket:
ID: %s
User ID: %s
Title: %s
Description: %s
`,
		input.TicketID,
		input.UserID,
		input.Title,
		input.Description,
	)
}
