package campaigns

import (
	"context"
	"testing"
	"time"
)

func TestFileRepository_Save(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	ctx := context.Background()

	campaign := &Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Graph:      ValidGraph,
	}

	err := repo.Save(ctx, campaign)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify by getting
	saved, err := repo.Get(ctx, "test_campaign")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if saved.Name != campaign.Name {
		t.Errorf("Save() name = %v, want %v", saved.Name, campaign.Name)
	}
}

func TestFileRepository_Get(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	ctx := context.Background()

	// Test get non-existent
	result, err := repo.Get(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if result != nil {
		t.Error("Get() should return nil for non-existent campaign")
	}

	// Save and get
	campaign := &Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit", "anticipation"},
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Graph:      ValidGraph,
	}
	repo.Save(ctx, campaign)

	result, err = repo.Get(ctx, "test_campaign")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if result == nil {
		t.Fatal("Get() should return campaign")
	}
	if len(result.Operations) != 2 {
		t.Errorf("Get() operations = %v, want 2", len(result.Operations))
	}
}

func TestFileRepository_List(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	ctx := context.Background()

	// Test empty list
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) != 0 {
		t.Errorf("List() count = %v, want 0", len(list))
	}

	// Add campaigns
	for _, name := range []string{"campaign1", "campaign2", "campaign3"} {
		repo.Save(ctx, &Campaign{
			Name:       name,
			Operations: []string{"credit"},
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
			Graph:      ValidGraph,
		})
	}

	list, err = repo.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(list) != 3 {
		t.Errorf("List() count = %v, want 3", len(list))
	}
}

func TestFileRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	ctx := context.Background()

	// Save campaign
	campaign := &Campaign{
		Name:       "test_campaign",
		Operations: []string{"credit"},
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Graph:      ValidGraph,
	}
	repo.Save(ctx, campaign)

	// Delete
	err := repo.Delete(ctx, "test_campaign")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deleted
	result, _ := repo.Get(ctx, "test_campaign")
	if result != nil {
		t.Error("Delete() campaign should not exist after deletion")
	}
}

func TestFileRepository_ListByOperation(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	ctx := context.Background()

	// Add campaigns with different operations
	campaigns := []*Campaign{
		{Name: "credit_only", Operations: []string{"credit"}, Graph: ValidGraph, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		{Name: "anticipation_only", Operations: []string{"anticipation"}, Graph: ValidGraph, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
		{Name: "both", Operations: []string{"credit", "anticipation"}, Graph: ValidGraph, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()},
	}
	for _, c := range campaigns {
		repo.Save(ctx, c)
	}

	// Test credit filter
	creditCampaigns, err := repo.ListByOperation(ctx, "credit")
	if err != nil {
		t.Fatalf("ListByOperation() error = %v", err)
	}
	if len(creditCampaigns) != 2 {
		t.Errorf("ListByOperation(credit) count = %v, want 2", len(creditCampaigns))
	}

	// Test anticipation filter
	anticipationCampaigns, err := repo.ListByOperation(ctx, "anticipation")
	if err != nil {
		t.Fatalf("ListByOperation() error = %v", err)
	}
	if len(anticipationCampaigns) != 2 {
		t.Errorf("ListByOperation(anticipation) count = %v, want 2", len(anticipationCampaigns))
	}
}


