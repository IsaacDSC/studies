package blocklist

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"rulesengine/internal/shared"
)

func TestService_Block(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	userID := uuid.New()
	ctx := context.Background()

	// Test block
	user, err := service.Block(ctx, userID)
	if err != nil {
		t.Fatalf("Block() error = %v", err)
	}
	if user.UserID != userID {
		t.Errorf("Block() user.UserID = %v, want %v", user.UserID, userID)
	}

	// Test blocking same user again
	_, err = service.Block(ctx, userID)
	if err != shared.ErrUserAlreadyBlocked {
		t.Errorf("Block() error = %v, want %v", err, shared.ErrUserAlreadyBlocked)
	}
}

func TestService_Unblock(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	userID := uuid.New()
	ctx := context.Background()

	// Block user first
	_, err := service.Block(ctx, userID)
	if err != nil {
		t.Fatalf("Block() error = %v", err)
	}

	// Test unblock
	err = service.Unblock(ctx, userID)
	if err != nil {
		t.Fatalf("Unblock() error = %v", err)
	}

	// Verify user is not blocked
	isBlocked, _ := service.IsBlocked(ctx, userID)
	if isBlocked {
		t.Error("Unblock() user should not be blocked after unblock")
	}
}

func TestService_Unblock_NotBlocked(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	userID := uuid.New()
	ctx := context.Background()

	// Test unblock non-blocked user
	err := service.Unblock(ctx, userID)
	if err != shared.ErrUserNotBlocked {
		t.Errorf("Unblock() error = %v, want %v", err, shared.ErrUserNotBlocked)
	}
}

func TestService_IsBlocked(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	service := NewService(repo)
	userID := uuid.New()
	ctx := context.Background()

	// Test not blocked
	isBlocked, err := service.IsBlocked(ctx, userID)
	if err != nil {
		t.Fatalf("IsBlocked() error = %v", err)
	}
	if isBlocked {
		t.Error("IsBlocked() should return false for non-blocked user")
	}

	// Block user
	_, err = service.Block(ctx, userID)
	if err != nil {
		t.Fatalf("Block() error = %v", err)
	}

	// Test blocked
	isBlocked, err = service.IsBlocked(ctx, userID)
	if err != nil {
		t.Fatalf("IsBlocked() error = %v", err)
	}
	if !isBlocked {
		t.Error("IsBlocked() should return true for blocked user")
	}
}


