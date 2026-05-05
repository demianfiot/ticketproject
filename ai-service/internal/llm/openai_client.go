package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/demianfiot/ticketproject/ai-service/internal/domain"
	"github.com/demianfiot/ticketproject/ai-service/internal/parser"
	"github.com/demianfiot/ticketproject/ai-service/internal/prompt"
)

type OpenAIClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
	prompt     *prompt.Builder
	parser     *parser.Parser
}

func NewOpenAIClient(
	apiKey string,
	model string,
	baseURL string,
	timeout time.Duration,
	promptBuilder *prompt.Builder,
	parser *parser.Parser,
) *OpenAIClient {
	return &OpenAIClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		prompt: promptBuilder,
		parser: parser,
	}
}

func (c *OpenAIClient) AnalyzeTicket(ctx context.Context, input domain.AnalyzeInput) (domain.AnalyzeResult, error) {
	userPrompt := c.prompt.BuildTicketAnalysisPrompt(input)

	reqBody := map[string]interface{}{
		"model": c.model,
		"input": []map[string]string{
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"temperature":       0.2,
		"max_output_tokens": 500,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("failed to marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("failed to create openai request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("failed to call openai: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("failed to read openai response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return domain.AnalyzeResult{}, fmt.Errorf("openai returned status %d: %s", resp.StatusCode, string(respBody))
	}

	text, err := extractResponseText(respBody)
	if err != nil {
		return domain.AnalyzeResult{}, err
	}

	result, err := c.parser.ParseAnalysis(text)
	if err != nil {
		return domain.AnalyzeResult{}, fmt.Errorf("failed to parse analysis: %w", err)
	}

	return result, nil
}

func extractResponseText(body []byte) (string, error) {
	var response struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal openai response: %w", err)
	}

	if response.OutputText != "" {
		return response.OutputText, nil
	}

	for _, out := range response.Output {
		for _, content := range out.Content {
			if content.Text != "" {
				return content.Text, nil
			}
		}
	}

	return "", fmt.Errorf("openai response text is empty")
}
