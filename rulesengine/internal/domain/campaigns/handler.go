package campaigns

import (
	"io"
	"net/http"
	"strings"

	"rulesengine/internal/shared"
)

// Handler handles HTTP requests for campaigns operations
type Handler struct {
	service *Service
}

// NewHandler creates a new campaigns handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Upsert handles PATCH /backoffice/campaigns/{campaignName}
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	campaignName := r.PathValue("campaignName")
	if campaignName == "" {
		shared.Error(w, shared.NewAppError("invalid_campaign", "campaign name is required", 400))
		return
	}

	// Get operations from query param
	operationsStr := r.URL.Query().Get("operations")
	if operationsStr == "" {
		shared.Error(w, shared.ErrMissingOperations)
		return
	}

	operations := strings.Split(operationsStr, ",")
	for i, op := range operations {
		operations[i] = strings.TrimSpace(op)
	}

	// Validate operations
	for _, op := range operations {
		if op != "credit" && op != "anticipation" {
			shared.Error(w, shared.ErrInvalidOperation)
			return
		}
	}

	// Read body as graph
	body, err := io.ReadAll(r.Body)
	if err != nil {
		shared.Error(w, shared.ErrInvalidJSON)
		return
	}

	campaign, err := h.service.Upsert(r.Context(), campaignName, operations, body)
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	shared.Success(w, UpsertResponse{
		Message:    "rules updated successfully",
		Campaign:   campaign.Name,
		Operations: campaign.Operations,
		UpdatedAt:  campaign.UpdatedAt,
	})
}

// List handles GET /backoffice/campaigns/list
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.service.List(r.Context())
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	summaries := make([]CampaignSummary, 0, len(campaigns))
	for _, c := range campaigns {
		summaries = append(summaries, c.ToSummary())
	}

	shared.Success(w, ListResponse{
		Campaigns: summaries,
		Total:     len(summaries),
	})
}

// Delete handles DELETE /backoffice/campaigns/{campaignName}/{operation}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	campaignName := r.PathValue("campaignName")
	operation := r.PathValue("operation")

	if campaignName == "" {
		shared.Error(w, shared.NewAppError("invalid_campaign", "campaign name is required", 400))
		return
	}

	if operation != "credit" && operation != "anticipation" {
		shared.Error(w, shared.ErrInvalidOperation)
		return
	}

	err := h.service.Delete(r.Context(), campaignName)
	if err != nil {
		shared.ErrorFromErr(w, err)
		return
	}

	shared.Success(w, DeleteResponse{
		Message:   "campaign deleted successfully",
		Campaign:  campaignName,
		Operation: operation,
	})
}


