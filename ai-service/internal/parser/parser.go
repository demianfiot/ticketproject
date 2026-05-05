package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseAnalysis(raw string) (domain.AnalyzeResult, error) {
	cleaned := cleanJSON(raw)

	var result domain.AnalyzeResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("invalid json from llm: %w, raw=%s", err, raw)
	}

	result.Category = normalizeCategory(result.Category)
	result.Priority = normalizePriority(result.Priority)

	if result.Summary == "" {
		result.Summary = "No summary provided"
	}

	if result.SuggestedReply == "" {
		result.SuggestedReply = "Support agent should review this ticket manually."
	}

	return result, nil
}

func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}

func normalizeCategory(category string) string {
	category = strings.ToLower(strings.TrimSpace(category))

	switch category {
	case "billing", "technical", "account", "refund", "other":
		return category
	default:
		return "other"
	}
}

func normalizePriority(priority string) string {
	priority = strings.ToLower(strings.TrimSpace(priority))

	switch priority {
	case "low", "medium", "high", "urgent":
		return priority
	default:
		return "medium"
	}
}
