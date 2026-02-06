package blocklist

import (
	"time"

	"github.com/google/uuid"
)

// BlockedUser represents a blocked user
type BlockedUser struct {
	UserID    uuid.UUID `json:"user_id"`
	BlockedAt time.Time `json:"blocked_at"`
}

// BlockRequest represents the request to block a user
type BlockRequest struct {
	UserID string `json:"user_id"`
}

// BlockResponse represents the response after blocking a user
type BlockResponse struct {
	Message   string    `json:"message"`
	UserID    string    `json:"user_id"`
	BlockedAt time.Time `json:"blocked_at"`
}

// UnblockResponse represents the response after unblocking a user
type UnblockResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}


