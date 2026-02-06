package blocklist

import (
	"net/http"

	"github.com/google/uuid"

	"rulesengine/internal/shared"
)

// Handler handles HTTP requests for blocklist operations
type Handler struct {
	service *Service
}

// NewHandler creates a new blocklist handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Block handles POST /backoffice/block-list
func (h *Handler) Block(w http.ResponseWriter, r *http.Request) {
	var req BlockRequest
	if err := shared.DecodeJSON(r, &req); err != nil {
		shared.Error(w, shared.ErrInvalidJSON)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		shared.Error(w, shared.ErrInvalidUUID)
		return
	}

	blockedUser, err := h.service.Block(r.Context(), userID)
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	shared.Created(w, BlockResponse{
		Message:   "user blocked successfully",
		UserID:    blockedUser.UserID.String(),
		BlockedAt: blockedUser.BlockedAt,
	})
}

// Unblock handles DELETE /backoffice/block-list/{user_id}
func (h *Handler) Unblock(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		shared.Error(w, shared.ErrInvalidUUID)
		return
	}

	err = h.service.Unblock(r.Context(), userID)
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	shared.Success(w, UnblockResponse{
		Message: "user unblocked successfully",
		UserID:  userID.String(),
	})
}


