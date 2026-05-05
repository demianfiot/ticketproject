package domain

type AnalyzeInput struct {
	TicketID    string
	UserID      string
	Title       string
	Description string
}

type AnalyzeResult struct {
	Summary        string  `json:"summary"`
	Category       string  `json:"category"`
	Priority       string  `json:"priority"`
	SuggestedReply string  `json:"suggested_reply"`
	Confidence     float64 `json:"confidence"`
}
