package operation

import (
	"time"

	"github.com/google/uuid"
)

// ValidOperations defines the allowed operations
var ValidOperations = map[string]bool{
	"credit":       true,
	"anticipation": true,
}

// EvaluationRequest represents the request to evaluate rules
type EvaluationRequest struct {
	RequestedAmount int    `json:"requested_amount"`
	UserID          string `json:"user_id"`
}

// EvaluationInput represents the validated input for evaluation
type EvaluationInput struct {
	RequestedAmount int
	UserID          uuid.UUID
}

// EvaluationResult represents the result of rule evaluation
type EvaluationResult struct {
	Status          string    `json:"status"` // approved, denied, blocked
	UserID          string    `json:"user_id"`
	RequestedAmount int       `json:"requested_amount"`
	Offers          []Offer   `json:"offers"`
	Reason          string    `json:"reason,omitempty"`
	Message         string    `json:"message,omitempty"`
	EvaluatedAt     time.Time `json:"evaluated_at"`
}

// Offer represents a credit/anticipation offer from a campaign
type Offer struct {
	Campaign     string  `json:"campaign"`
	Approved     bool    `json:"approved"`
	Rate         float64 `json:"rate"`
	Installments []int   `json:"installments"`
	Value        int     `json:"value"`
}

// NewBlockedResult creates a blocked result
func NewBlockedResult(userID string, requestedAmount int) *EvaluationResult {
	return &EvaluationResult{
		Status:          "blocked",
		UserID:          userID,
		RequestedAmount: requestedAmount,
		Reason:          "user_blocked",
		Message:         "User is blocked for credit operations",
		Offers:          []Offer{},
		EvaluatedAt:     time.Now().UTC(),
	}
}

// NewDeniedResult creates a denied result
func NewDeniedResult(userID string, requestedAmount int, reason, message string) *EvaluationResult {
	return &EvaluationResult{
		Status:          "denied",
		UserID:          userID,
		RequestedAmount: requestedAmount,
		Reason:          reason,
		Message:         message,
		Offers:          []Offer{},
		EvaluatedAt:     time.Now().UTC(),
	}
}

// NewApprovedResult creates an approved result with offers
func NewApprovedResult(userID string, requestedAmount int, offers []Offer) *EvaluationResult {
	return &EvaluationResult{
		Status:          "approved",
		UserID:          userID,
		RequestedAmount: requestedAmount,
		Offers:          offers,
		EvaluatedAt:     time.Now().UTC(),
	}
}


