package operation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"rulesengine/internal/domain/blocklist"
	"rulesengine/internal/domain/campaigns"
	"rulesengine/internal/shared"
)

func setupTestEngine(t *testing.T) (*Engine, *blocklist.Service, *campaigns.FileRepository) {
	blocklistDir := t.TempDir()
	campaignsDir := t.TempDir()

	blocklistRepo := blocklist.NewFileRepository(blocklistDir)
	blocklistSvc := blocklist.NewService(blocklistRepo)

	campaignsRepo := campaigns.NewFileRepository(campaignsDir)

	engine := NewEngine(campaignsRepo, blocklistSvc)
	return engine, blocklistSvc, campaignsRepo
}

func TestEngine_Evaluate_BlockedUser(t *testing.T) {
	engine, blocklistSvc, _ := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	// Block user
	blocklistSvc.Block(ctx, userID)

	input := EvaluationInput{
		RequestedAmount: 100000,
		UserID:          userID,
	}

	result, err := engine.Evaluate(ctx, "credit", input)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.Status != "blocked" {
		t.Errorf("Evaluate() status = %v, want blocked", result.Status)
	}
	if result.Reason != "user_blocked" {
		t.Errorf("Evaluate() reason = %v, want user_blocked", result.Reason)
	}
}

func TestEngine_Evaluate_NoCampaigns(t *testing.T) {
	engine, _, _ := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	input := EvaluationInput{
		RequestedAmount: 100000,
		UserID:          userID,
	}

	_, err := engine.Evaluate(ctx, "credit", input)
	if err != shared.ErrNoCampaigns {
		t.Errorf("Evaluate() error = %v, want %v", err, shared.ErrNoCampaigns)
	}
}

func TestEngine_Evaluate_Approved(t *testing.T) {
	engine, _, campaignsRepo := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	// Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	campaignsRepo.Save(ctx, campaign)

	input := EvaluationInput{
		RequestedAmount: 50000, // 500 dollars, should be approved
		UserID:          userID,
	}

	result, err := engine.Evaluate(ctx, "credit", input)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.Status != "approved" {
		t.Errorf("Evaluate() status = %v, want approved", result.Status)
	}
	if len(result.Offers) != 1 {
		t.Fatalf("Evaluate() offers count = %v, want 1", len(result.Offers))
	}
	if !result.Offers[0].Approved {
		t.Error("Evaluate() offer should be approved")
	}
	if result.Offers[0].Rate != 1.99 {
		t.Errorf("Evaluate() rate = %v, want 1.99", result.Offers[0].Rate)
	}
}

func TestEngine_Evaluate_Denied(t *testing.T) {
	engine, _, campaignsRepo := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	// Create campaign
	campaign := &campaigns.Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	campaignsRepo.Save(ctx, campaign)

	input := EvaluationInput{
		RequestedAmount: 200000, // 2000 dollars, should be denied
		UserID:          userID,
	}

	result, err := engine.Evaluate(ctx, "credit", input)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.Status != "denied" {
		t.Errorf("Evaluate() status = %v, want denied", result.Status)
	}
	if result.Reason != "credit_limit_exceeded" {
		t.Errorf("Evaluate() reason = %v, want credit_limit_exceeded", result.Reason)
	}
}

func TestEngine_Evaluate_MultipleCampaigns(t *testing.T) {
	engine, _, campaignsRepo := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	// Create multiple campaigns
	for _, name := range []string{"campaign1", "campaign2", "campaign3"} {
		campaign := &campaigns.Campaign{
			Name:       name,
			Operations: []string{"credit"},
			Graph:      ValidCreditGraph,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		campaignsRepo.Save(ctx, campaign)
	}

	input := EvaluationInput{
		RequestedAmount: 50000,
		UserID:          userID,
	}

	result, err := engine.Evaluate(ctx, "credit", input)
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	if result.Status != "approved" {
		t.Errorf("Evaluate() status = %v, want approved", result.Status)
	}
	if len(result.Offers) != 3 {
		t.Errorf("Evaluate() offers count = %v, want 3", len(result.Offers))
	}
}

func TestEngine_Evaluate_FiltersByOperation(t *testing.T) {
	engine, _, campaignsRepo := setupTestEngine(t)
	ctx := context.Background()
	userID := uuid.New()

	// Create campaigns with different operations
	creditCampaign := &campaigns.Campaign{
		Name:       "credit_only",
		Operations: []string{"credit"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	anticipationCampaign := &campaigns.Campaign{
		Name:       "anticipation_only",
		Operations: []string{"anticipation"},
		Graph:      ValidCreditGraph,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	campaignsRepo.Save(ctx, creditCampaign)
	campaignsRepo.Save(ctx, anticipationCampaign)

	input := EvaluationInput{
		RequestedAmount: 50000,
		UserID:          userID,
	}

	// Test credit operation
	result, err := engine.Evaluate(ctx, "credit", input)
	if err != nil {
		t.Fatalf("Evaluate(credit) error = %v", err)
	}
	if len(result.Offers) != 1 {
		t.Errorf("Evaluate(credit) offers count = %v, want 1", len(result.Offers))
	}
	if result.Offers[0].Campaign != "credit_only" {
		t.Errorf("Evaluate(credit) campaign = %v, want credit_only", result.Offers[0].Campaign)
	}

	// Test anticipation operation
	result, err = engine.Evaluate(ctx, "anticipation", input)
	if err != nil {
		t.Fatalf("Evaluate(anticipation) error = %v", err)
	}
	if len(result.Offers) != 1 {
		t.Errorf("Evaluate(anticipation) offers count = %v, want 1", len(result.Offers))
	}
	if result.Offers[0].Campaign != "anticipation_only" {
		t.Errorf("Evaluate(anticipation) campaign = %v, want anticipation_only", result.Offers[0].Campaign)
	}
}


