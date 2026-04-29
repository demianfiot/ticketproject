package http

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type Ticket struct {
	ID          int    `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`

	AISummary        string `json:"ai_summary"`
	AICategory       string `json:"ai_category"`
	AIPriority       string `json:"ai_priority"`
	AISuggestedReply string `json:"ai_suggested_reply"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}