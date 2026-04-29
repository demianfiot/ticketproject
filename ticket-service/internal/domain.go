package domain

import (
	"time"
)

type Ticket struct {
	ID               int       `db:"id"`
	UserID           string    `db:"user_id"`
	Title            string    `db:"title"`
	Description      string    `db:"description"`
	Status           string    `db:"status"`
	AISummary        string    `db:"ai_summary"`
	AICategory       string    `db:"ai_category"`
	AIPriority       string    `db:"ai_priority"`
	AISuggestedReply string    `db:"ai_suggested_reply"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
