package campaigns

import (
	"context"
	"testing"

	"rulesengine/internal/shared"
)

func TestService_Upsert(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	// Test create
	campaign, err := service.Upsert(ctx, "test_campaign", []string{"credit"}, ValidGraph)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if campaign.Name != "test_campaign" {
		t.Errorf("Upsert() name = %v, want test_campaign", campaign.Name)
	}
	if len(campaign.Operations) != 1 {
		t.Errorf("Upsert() operations count = %v, want 1", len(campaign.Operations))
	}

	// Test update
	campaign, err = service.Upsert(ctx, "test_campaign", []string{"credit", "anticipation"}, ValidGraph)
	if err != nil {
		t.Fatalf("Upsert() update error = %v", err)
	}
	if len(campaign.Operations) != 2 {
		t.Errorf("Upsert() updated operations count = %v, want 2", len(campaign.Operations))
	}
}

func TestService_Upsert_InvalidGraph(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	_, err := service.Upsert(ctx, "test_campaign", []string{"credit"}, InvalidGraph)
	if err != shared.ErrInvalidGraph {
		t.Errorf("Upsert() error = %v, want %v", err, shared.ErrInvalidGraph)
	}
}

func TestService_Get(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	// Test not found
	_, err := service.Get(ctx, "nonexistent")
	if err != shared.ErrCampaignNotFound {
		t.Errorf("Get() error = %v, want %v", err, shared.ErrCampaignNotFound)
	}

	// Create and get
	service.Upsert(ctx, "test_campaign", []string{"credit"}, ValidGraph)

	campaign, err := service.Get(ctx, "test_campaign")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if campaign.Name != "test_campaign" {
		t.Errorf("Get() name = %v, want test_campaign", campaign.Name)
	}
}

func TestService_List(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	// Test empty list
	campaigns, err := service.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(campaigns) != 0 {
		t.Errorf("List() count = %v, want 0", len(campaigns))
	}

	// Add campaigns
	service.Upsert(ctx, "campaign1", []string{"credit"}, ValidGraph)
	service.Upsert(ctx, "campaign2", []string{"anticipation"}, ValidGraph)

	campaigns, err = service.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(campaigns) != 2 {
		t.Errorf("List() count = %v, want 2", len(campaigns))
	}
}

func TestService_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	// Test delete non-existent
	err := service.Delete(ctx, "nonexistent")
	if err != shared.ErrCampaignNotFound {
		t.Errorf("Delete() error = %v, want %v", err, shared.ErrCampaignNotFound)
	}

	// Create and delete
	service.Upsert(ctx, "test_campaign", []string{"credit"}, ValidGraph)

	err = service.Delete(ctx, "test_campaign")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deleted
	_, err = service.Get(ctx, "test_campaign")
	if err != shared.ErrCampaignNotFound {
		t.Errorf("Delete() campaign should not exist after deletion")
	}
}

func TestService_ListByOperation(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	ctx := context.Background()

	// Add campaigns
	service.Upsert(ctx, "credit_only", []string{"credit"}, ValidGraph)
	service.Upsert(ctx, "both", []string{"credit", "anticipation"}, ValidGraph)

	campaigns, err := service.ListByOperation(ctx, "credit")
	if err != nil {
		t.Fatalf("ListByOperation() error = %v", err)
	}
	if len(campaigns) != 2 {
		t.Errorf("ListByOperation() count = %v, want 2", len(campaigns))
	}

	campaigns, err = service.ListByOperation(ctx, "anticipation")
	if err != nil {
		t.Fatalf("ListByOperation() error = %v", err)
	}
	if len(campaigns) != 1 {
		t.Errorf("ListByOperation() count = %v, want 1", len(campaigns))
	}
}


