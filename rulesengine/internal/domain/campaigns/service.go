package campaigns

import (
	"context"
	"encoding/json"
	"time"

	zen "github.com/gorules/zen-go"

	"rulesengine/internal/shared"
)

// Service handles campaigns business logic
type Service struct {
	repo   Repository
	engine zen.Engine
}

// NewService creates a new campaigns service
func NewService(repo Repository) *Service {
	return &Service{
		repo:   repo,
		engine: zen.NewEngine(zen.EngineConfig{}),
	}
}

// Upsert creates or updates a campaign
func (s *Service) Upsert(ctx context.Context, name string, operations []string, graph json.RawMessage) (*Campaign, error) {
	// Validate graph by trying to create a decision
	_, err := s.engine.CreateDecision(graph)
	if err != nil {
		return nil, shared.ErrInvalidGraph
	}

	// Get existing campaign or create new
	existing, err := s.repo.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	campaign := &Campaign{
		Name:       name,
		Operations: operations,
		Graph:      graph,
		UpdatedAt:  now,
	}

	if existing != nil {
		campaign.CreatedAt = existing.CreatedAt
	} else {
		campaign.CreatedAt = now
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

// Get retrieves a campaign by name
func (s *Service) Get(ctx context.Context, name string) (*Campaign, error) {
	campaign, err := s.repo.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, shared.ErrCampaignNotFound
	}
	return campaign, nil
}

// List retrieves all campaigns
func (s *Service) List(ctx context.Context) ([]*Campaign, error) {
	return s.repo.List(ctx)
}

// Delete removes a campaign
func (s *Service) Delete(ctx context.Context, name string) error {
	// Check if campaign exists
	existing, err := s.repo.Get(ctx, name)
	if err != nil {
		return err
	}
	if existing == nil {
		return shared.ErrCampaignNotFound
	}

	return s.repo.Delete(ctx, name)
}

// ListByOperation retrieves all campaigns for an operation
func (s *Service) ListByOperation(ctx context.Context, operation string) ([]*Campaign, error) {
	return s.repo.ListByOperation(ctx, operation)
}


