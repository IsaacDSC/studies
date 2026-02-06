package campaigns

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Repository defines the interface for campaigns storage
type Repository interface {
	Save(ctx context.Context, campaign *Campaign) error
	Get(ctx context.Context, name string) (*Campaign, error)
	List(ctx context.Context) ([]*Campaign, error)
	Delete(ctx context.Context, name string) error
	ListByOperation(ctx context.Context, operation string) ([]*Campaign, error)
}

// FileRepository implements Repository using file system storage
type FileRepository struct {
	basePath string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(basePath string) *FileRepository {
	os.MkdirAll(basePath, 0755)
	return &FileRepository{basePath: basePath}
}

func (r *FileRepository) filePath(name string) string {
	// Sanitize name to prevent path traversal
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, "..", "_")
	return filepath.Join(r.basePath, safeName+".json")
}

// Save saves a campaign to file
func (r *FileRepository) Save(ctx context.Context, campaign *Campaign) error {
	data, err := json.MarshalIndent(campaign, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath(campaign.Name), data, 0644)
}

// Get retrieves a campaign by name
func (r *FileRepository) Get(ctx context.Context, name string) (*Campaign, error) {
	data, err := os.ReadFile(r.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var campaign Campaign
	if err := json.Unmarshal(data, &campaign); err != nil {
		return nil, err
	}
	return &campaign, nil
}

// List retrieves all campaigns
func (r *FileRepository) List(ctx context.Context) ([]*Campaign, error) {
	entries, err := os.ReadDir(r.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Campaign{}, nil
		}
		return nil, err
	}

	var campaigns []*Campaign
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(r.basePath, entry.Name()))
		if err != nil {
			continue
		}

		var campaign Campaign
		if err := json.Unmarshal(data, &campaign); err != nil {
			continue
		}
		campaigns = append(campaigns, &campaign)
	}
	return campaigns, nil
}

// Delete removes a campaign file
func (r *FileRepository) Delete(ctx context.Context, name string) error {
	return os.Remove(r.filePath(name))
}

// ListByOperation retrieves all campaigns that have the specified operation
func (r *FileRepository) ListByOperation(ctx context.Context, operation string) ([]*Campaign, error) {
	allCampaigns, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*Campaign
	for _, c := range allCampaigns {
		if c.HasOperation(operation) {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}


