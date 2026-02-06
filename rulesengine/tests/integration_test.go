package tests

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
	"rulesengine/internal/domain/operation"
	"rulesengine/internal/middleware"
)

// ValidCreditGraph is a valid gorules decision graph for testing
var ValidCreditGraph = json.RawMessage(`{
  "contentType": "application/vnd.gorules.decision",
  "nodes": [
    {"id": "input", "type": "inputNode", "name": "Request", "position": {"x": 0, "y": 0}},
    {"id": "output", "type": "outputNode", "name": "Response", "position": {"x": 0, "y": 0}},
    {
      "id": "credit_check",
      "type": "decisionTableNode",
      "name": "Credit Check",
      "position": {"x": 0, "y": 0},
      "content": {
        "hitPolicy": "first",
        "inputs": [
          {"id": "amount", "name": "Amount", "field": "requested_amount"}
        ],
        "outputs": [
          {"id": "approved", "name": "Approved", "field": "approved"},
          {"id": "rate", "name": "Rate", "field": "rate"},
          {"id": "installments", "name": "Installments", "field": "installments"},
          {"id": "value", "name": "Value", "field": "value"}
        ],
        "rules": [
          {"id": "r1", "amount": "<= 100000", "approved": "true", "rate": "1.99", "installments": "[6, 12, 24]", "value": "requested_amount"},
          {"id": "r2", "amount": "<= 500000", "approved": "true", "rate": "2.99", "installments": "[6, 12]", "value": "requested_amount"},
          {"id": "r3", "amount": "> 500000", "approved": "false", "rate": "0", "installments": "[]", "value": "0"}
        ]
      }
    }
  ],
  "edges": [
    {"id": "e1", "sourceId": "input", "targetId": "credit_check"},
    {"id": "e2", "sourceId": "credit_check", "targetId": "output"}
  ]
}`)

// TestServer wraps all the server components for integration testing
type TestServer struct {
	Handler          http.Handler
	BlocklistSvc     *blocklist.Service
	CampaignsSvc     *campaigns.Service
	CampaignsRepo    *campaigns.FileRepository
	BlocklistHandler *blocklist.Handler
	CampaignsHandler *campaigns.Handler
	OperationHandler *operation.Handler
}

func setupTestServer(t *testing.T) *TestServer {
	blocklistDir := t.TempDir()
	campaignsDir := t.TempDir()

	// Initialize repositories
	blocklistRepo := blocklist.NewFileRepository(blocklistDir)
	campaignsRepo := campaigns.NewFileRepository(campaignsDir)

	// Initialize services
	blocklistSvc := blocklist.NewService(blocklistRepo)
	campaignsSvc := campaigns.NewService(campaignsRepo)

	// Initialize engine
	operationEngine := operation.NewEngine(campaignsRepo, blocklistSvc)

	// Initialize handlers
	blocklistHandler := blocklist.NewHandler(blocklistSvc)
	campaignsHandler := campaigns.NewHandler(campaignsSvc)
	operationHandler := operation.NewHandler(operationEngine)

	// Setup routes
	mux := http.NewServeMux()

	// Blocklist routes
	mux.HandleFunc("POST /backoffice/block-list", blocklistHandler.Block)
	mux.HandleFunc("DELETE /backoffice/block-list/{user_id}", blocklistHandler.Unblock)

	// Campaigns routes
	mux.HandleFunc("PATCH /backoffice/campaigns/{campaignName}", campaignsHandler.Upsert)
	mux.HandleFunc("GET /backoffice/campaigns/list", campaignsHandler.List)
	mux.HandleFunc("DELETE /backoffice/campaigns/{campaignName}/{operation}", campaignsHandler.Delete)

	// Operation routes
	mux.HandleFunc("POST /rule-engine/rule/{operation}", operationHandler.Evaluate)

	// Apply middleware
	handler := middleware.Recovery(middleware.Logging(mux))

	return &TestServer{
		Handler:          handler,
		BlocklistSvc:     blocklistSvc,
		CampaignsSvc:     campaignsSvc,
		CampaignsRepo:    campaignsRepo,
		BlocklistHandler: blocklistHandler,
		CampaignsHandler: campaignsHandler,
		OperationHandler: operationHandler,
	}
}

// Test: Block user -> Request credit -> Expect blocked status
func TestIntegration_BlockedUser(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create campaign first
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	ts.CampaignsRepo.Save(nil, campaign)

	// 2. Block user
	blockBody, _ := json.Marshal(map[string]string{"user_id": userID.String()})
	blockReq := httptest.NewRequest(http.MethodPost, "/backoffice/block-list", bytes.NewReader(blockBody))
	blockReq.Header.Set("Content-Type", "application/json")
	blockRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(blockRR, blockReq)

	if blockRR.Code != http.StatusCreated {
		t.Fatalf("Block user: status = %d, want %d", blockRR.Code, http.StatusCreated)
	}

	// 3. Request credit - should be blocked
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 50000,
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusOK {
		t.Fatalf("Request credit: status = %d, want %d", creditRR.Code, http.StatusOK)
	}

	var result map[string]any
	json.NewDecoder(creditRR.Body).Decode(&result)
	if result["status"] != "blocked" {
		t.Errorf("Request credit: status = %v, want blocked", result["status"])
	}
}

// Test: Create campaign -> Request credit -> Expect offers
func TestIntegration_CreateCampaignAndRequestCredit(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create campaign via API
	campaignReq := httptest.NewRequest(http.MethodPatch, "/backoffice/campaigns/black_friday?operations=credit", bytes.NewReader(ValidCreditGraph))
	campaignReq.Header.Set("Content-Type", "application/json")
	campaignRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(campaignRR, campaignReq)

	if campaignRR.Code != http.StatusOK {
		t.Fatalf("Create campaign: status = %d, want %d, body: %s", campaignRR.Code, http.StatusOK, campaignRR.Body.String())
	}

	// 2. Request credit
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 50000,
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusOK {
		t.Fatalf("Request credit: status = %d, want %d, body: %s", creditRR.Code, http.StatusOK, creditRR.Body.String())
	}

	var result map[string]any
	json.NewDecoder(creditRR.Body).Decode(&result)
	if result["status"] != "approved" {
		t.Errorf("Request credit: status = %v, want approved", result["status"])
	}

	offers, ok := result["offers"].([]any)
	if !ok || len(offers) == 0 {
		t.Error("Request credit: should have offers")
	}
}

// Test: Multiple campaigns -> Parallel execution -> Multiple offers
func TestIntegration_MultipleCampaignsParallelExecution(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create multiple campaigns
	campaignNames := []string{"campaign1", "campaign2", "campaign3"}
	for _, name := range campaignNames {
		campaign := &campaigns.Campaign{
			Name:       name,
			Operations: []string{"credit"},
			Graph:      ValidCreditGraph,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		ts.CampaignsRepo.Save(nil, campaign)
	}

	// 2. Request credit
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 50000,
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusOK {
		t.Fatalf("Request credit: status = %d, want %d", creditRR.Code, http.StatusOK)
	}

	var result map[string]any
	json.NewDecoder(creditRR.Body).Decode(&result)

	offers, ok := result["offers"].([]any)
	if !ok {
		t.Fatal("Request credit: offers should be an array")
	}
	if len(offers) != 3 {
		t.Errorf("Request credit: offers count = %d, want 3", len(offers))
	}
}

// Test: Delete campaign -> Request credit -> No campaigns error
func TestIntegration_DeleteCampaignAndRequestCredit(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create campaign
	campaign := &campaigns.Campaign{
		Name:       "to_delete",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	ts.CampaignsRepo.Save(nil, campaign)

	// 2. Delete campaign via API
	deleteReq := httptest.NewRequest(http.MethodDelete, "/backoffice/campaigns/to_delete/credit", nil)
	deleteRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusOK {
		t.Fatalf("Delete campaign: status = %d, want %d", deleteRR.Code, http.StatusOK)
	}

	// 3. Request credit - should fail with no campaigns
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 50000,
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusUnprocessableEntity {
		t.Errorf("Request credit: status = %d, want %d", creditRR.Code, http.StatusUnprocessableEntity)
	}
}

// Test: Unblock user -> Request credit -> Expect approved
func TestIntegration_UnblockUserAndRequestCredit(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	ts.CampaignsRepo.Save(nil, campaign)

	// 2. Block user
	ts.BlocklistSvc.Block(nil, userID)

	// 3. Unblock user via API
	unblockReq := httptest.NewRequest(http.MethodDelete, "/backoffice/block-list/"+userID.String(), nil)
	unblockRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(unblockRR, unblockReq)

	if unblockRR.Code != http.StatusOK {
		t.Fatalf("Unblock user: status = %d, want %d", unblockRR.Code, http.StatusOK)
	}

	// 4. Request credit - should be approved
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 50000,
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusOK {
		t.Fatalf("Request credit: status = %d, want %d", creditRR.Code, http.StatusOK)
	}

	var result map[string]any
	json.NewDecoder(creditRR.Body).Decode(&result)
	if result["status"] != "approved" {
		t.Errorf("Request credit: status = %v, want approved", result["status"])
	}
}

// Test: List campaigns
func TestIntegration_ListCampaigns(t *testing.T) {
	ts := setupTestServer(t)

	// 1. Create campaigns
	for _, name := range []string{"campaign1", "campaign2"} {
		campaign := &campaigns.Campaign{
			Name:       name,
			Operations: []string{"credit"},
			Graph:      ValidCreditGraph,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		ts.CampaignsRepo.Save(nil, campaign)
	}

	// 2. List campaigns
	listReq := httptest.NewRequest(http.MethodGet, "/backoffice/campaigns/list", nil)
	listRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(listRR, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("List campaigns: status = %d, want %d", listRR.Code, http.StatusOK)
	}

	var result map[string]any
	json.NewDecoder(listRR.Body).Decode(&result)

	total, ok := result["total"].(float64)
	if !ok || int(total) != 2 {
		t.Errorf("List campaigns: total = %v, want 2", result["total"])
	}
}

// Test: Credit denied when amount exceeds limit
func TestIntegration_CreditDenied(t *testing.T) {
	ts := setupTestServer(t)
	userID := uuid.New()

	// 1. Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	ts.CampaignsRepo.Save(nil, campaign)

	// 2. Request credit with high amount - should be denied
	creditBody, _ := json.Marshal(map[string]any{
		"requested_amount": 1000000, // Over 500000 limit
		"user_id":          userID.String(),
	})
	creditReq := httptest.NewRequest(http.MethodPost, "/rule-engine/rule/credit", bytes.NewReader(creditBody))
	creditReq.Header.Set("Content-Type", "application/json")
	creditRR := httptest.NewRecorder()
	ts.Handler.ServeHTTP(creditRR, creditReq)

	if creditRR.Code != http.StatusOK {
		t.Fatalf("Request credit: status = %d, want %d", creditRR.Code, http.StatusOK)
	}

	var result map[string]any
	json.NewDecoder(creditRR.Body).Decode(&result)
	if result["status"] != "denied" {
		t.Errorf("Request credit: status = %v, want denied", result["status"])
	}
	if result["reason"] != "credit_limit_exceeded" {
		t.Errorf("Request credit: reason = %v, want credit_limit_exceeded", result["reason"])
	}
}


