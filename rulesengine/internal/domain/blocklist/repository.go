package blocklist

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for blocklist storage
type Repository interface {
	Save(ctx context.Context, userID uuid.UUID) (*BlockedUser, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	Exists(ctx context.Context, userID uuid.UUID) (bool, error)
}

// FileRepository implements Repository using file system storage
type FileRepository struct {
	basePath string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(basePath string) *FileRepository {
	// Ensure directory exists
	os.MkdirAll(basePath, 0755)
	return &FileRepository{basePath: basePath}
}

func (r *FileRepository) filePath(userID uuid.UUID) string {
	return filepath.Join(r.basePath, userID.String()+".json")
}

// Save saves a blocked user to file
func (r *FileRepository) Save(ctx context.Context, userID uuid.UUID) (*BlockedUser, error) {
	user := &BlockedUser{
		UserID:    userID,
		BlockedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(r.filePath(userID), data, 0644)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a blocked user file
func (r *FileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return os.Remove(r.filePath(userID))
}

// Exists checks if a user is blocked
func (r *FileRepository) Exists(ctx context.Context, userID uuid.UUID) (bool, error) {
	_, err := os.Stat(r.filePath(userID))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}


