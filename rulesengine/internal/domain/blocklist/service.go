package blocklist

import (
	"context"

	"github.com/google/uuid"

	"rulesengine/internal/shared"
)

// Service handles blocklist business logic
type Service struct {
	repo Repository
}

// NewService creates a new blocklist service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Block blocks a user
func (s *Service) Block(ctx context.Context, userID uuid.UUID) (*BlockedUser, error) {
	// Check if user is already blocked
	exists, err := s.repo.Exists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, shared.ErrUserAlreadyBlocked
	}

	return s.repo.Save(ctx, userID)
}

// Unblock removes a user from the blocklist
func (s *Service) Unblock(ctx context.Context, userID uuid.UUID) error {
	// Check if user is blocked
	exists, err := s.repo.Exists(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return shared.ErrUserNotBlocked
	}

	return s.repo.Delete(ctx, userID)
}

// IsBlocked checks if a user is blocked
func (s *Service) IsBlocked(ctx context.Context, userID uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, userID)
}


