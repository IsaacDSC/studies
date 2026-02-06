package blocklist

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestFileRepository_Save(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	userID := uuid.New()
	ctx := context.Background()

	// Test Save
	user, err := repo.Save(ctx, userID)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if user.UserID != userID {
		t.Errorf("Save() user.UserID = %v, want %v", user.UserID, userID)
	}
	if user.BlockedAt.IsZero() {
		t.Error("Save() user.BlockedAt should not be zero")
	}

	// Verify file exists
	filePath := filepath.Join(tmpDir, userID.String()+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Save() file should exist")
	}
}

func TestFileRepository_Exists(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	userID := uuid.New()
	ctx := context.Background()

	// Test not exists
	exists, err := repo.Exists(ctx, userID)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() should return false for non-existent user")
	}

	// Save user
	_, err = repo.Save(ctx, userID)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Test exists
	exists, err = repo.Exists(ctx, userID)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() should return true for existing user")
	}
}

func TestFileRepository_Delete(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	userID := uuid.New()
	ctx := context.Background()

	// Save user first
	_, err := repo.Save(ctx, userID)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Delete user
	err = repo.Delete(ctx, userID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify file is gone
	exists, _ := repo.Exists(ctx, userID)
	if exists {
		t.Error("Delete() user should not exist after deletion")
	}
}

func TestFileRepository_Delete_NotExists(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	repo := NewFileRepository(tmpDir)
	userID := uuid.New()
	ctx := context.Background()

	// Delete non-existent user
	err := repo.Delete(ctx, userID)
	if err == nil {
		t.Error("Delete() should return error for non-existent user")
	}
}


