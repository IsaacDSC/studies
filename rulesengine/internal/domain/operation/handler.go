package operation

import (
	"net/http"

	"github.com/google/uuid"

	"rulesengine/internal/shared"
)

// Handler handles HTTP requests for operation evaluation
type Handler struct {
	engine *Engine
}

// NewHandler creates a new operation handler
func NewHandler(engine *Engine) *Handler {
	return &Handler{engine: engine}
}

// Evaluate handles POST /rule-engine/rule/{operation}
func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	operation := r.PathValue("operation")

	// Validate operation
	if !ValidOperations[operation] {
		shared.Error(w, shared.ErrInvalidOperation)
		return
	}

	// Parse request body
	var req EvaluationRequest
	if err := shared.DecodeJSON(r, &req); err != nil {
		shared.Error(w, shared.ErrInvalidJSON)
		return
	}

	// Validate user_id
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		shared.Error(w, shared.ErrInvalidUUID)
		return
	}

	// Validate requested_amount
	if req.RequestedAmount <= 0 {
		shared.Error(w, shared.ErrInvalidAmount)
		return
	}

	input := EvaluationInput{
		RequestedAmount: req.RequestedAmount,
		UserID:          userID,
	}

	// Evaluate rules
	result, err := h.engine.Evaluate(r.Context(), operation, input)
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	shared.Success(w, result)
}


