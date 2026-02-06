package operation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"rulesengine/internal/domain/blocklist"
	"rulesengine/internal/domain/campaigns"
)

func setupTestHandler(t *testing.T) (*Handler, *blocklist.Service, *campaigns.FileRepository) {
	blocklistDir := t.TempDir()
	campaignsDir := t.TempDir()

	blocklistRepo := blocklist.NewFileRepository(blocklistDir)
	blocklistSvc := blocklist.NewService(blocklistRepo)

	campaignsRepo := campaigns.NewFileRepository(campaignsDir)

	engine := NewEngine(campaignsRepo, blocklistSvc)
	handler := NewHandler(engine)

	return handler, blocklistSvc, campaignsRepo
}

func TestHandler_Evaluate(t *testing.T) {
	handler, _, campaignsRepo := setupTestHandler(t)
	ctx := campaignsRepo
	userID := uuid.New()

	// Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	ctx.Save(nil, campaign)

	// Create request
	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 50000,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Evaluate() status = %d, want %d, body: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var result EvaluationResult
	json.NewDecoder(rr.Body).Decode(&result)
	if result.Status != "approved" {
		t.Errorf("Evaluate() result.Status = %v, want approved", result.Status)
	}
}

func TestHandler_Evaluate_InvalidOperation(t *testing.T) {
	handler, _, _ := setupTestHandler(t)
	userID := uuid.New()

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 50000,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/invalid", bytes.NewReader(body))
	req.SetPathValue("operation", "invalid")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Evaluate_InvalidUUID(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 50000,
		UserID:          "invalid-uuid",
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Evaluate_InvalidAmount(t *testing.T) {
	handler, _, _ := setupTestHandler(t)
	userID := uuid.New()

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: -100,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Evaluate_ZeroAmount(t *testing.T) {
	handler, _, _ := setupTestHandler(t)
	userID := uuid.New()

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 0,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Evaluate_InvalidJSON(t *testing.T) {
	handler, _, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader([]byte("invalid json")))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestHandler_Evaluate_BlockedUser(t *testing.T) {
	handler, blocklistSvc, campaignsRepo := setupTestHandler(t)
	userID := uuid.New()

	// Block user
	blocklistSvc.Block(nil, userID)

	// Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	campaignsRepo.Save(nil, campaign)

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 50000,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusOK)
	}

	var result EvaluationResult
	json.NewDecoder(rr.Body).Decode(&result)
	if result.Status != "blocked" {
		t.Errorf("Evaluate() result.Status = %v, want blocked", result.Status)
	}
}

func TestHandler_Evaluate_NoCampaigns(t *testing.T) {
	handler, _, _ := setupTestHandler(t)
	userID := uuid.New()

	body, _ := json.Marshal(EvaluationRequest{
		RequestedAmount: 50000,
		UserID:          userID.String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(body))
	req.SetPathValue("operation", "credit")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Evaluate(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Evaluate() status = %d, want %d", rr.Code, http.StatusUnprocessableEntity)
	}
}


