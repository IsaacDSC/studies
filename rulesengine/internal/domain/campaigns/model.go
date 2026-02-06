package campaigns

import (
	"encoding/json"
	"time"
)

// Campaign represents a campaign with its rules graph
type Campaign struct {
	Name       string          `json:"campaign"`
	Operations []string        `json:"operations"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Graph      json.RawMessage `json:"graph"`
}

// CampaignSummary represents a campaign without the full graph
type CampaignSummary struct {
	Name       string    `json:"name"`
	Operations []string  `json:"operations"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UpsertRequest represents the request to create/update a campaign
type UpsertRequest struct {
	Graph json.RawMessage `json:"-"` // The entire body is the graph
}

// UpsertResponse represents the response after creating/updating a campaign
type UpsertResponse struct {
	Message    string    `json:"message"`
	Campaign   string    `json:"campaign"`
	Operations []string  `json:"operations"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ListResponse represents the response for listing campaigns
type ListResponse struct {
	Campaigns []CampaignSummary `json:"campaigns"`
	Total     int               `json:"total"`
}

// DeleteResponse represents the response after deleting a campaign
type DeleteResponse struct {
	Message   string `json:"message"`
	Campaign  string `json:"campaign"`
	Operation string `json:"operation"`
}

// HasOperation checks if the campaign has a specific operation
func (c *Campaign) HasOperation(operation string) bool {
	for _, op := range c.Operations {
		if op == operation {
			return true
		}
	}
	return false
}

// ToSummary converts a Campaign to CampaignSummary
func (c *Campaign) ToSummary() CampaignSummary {
	return CampaignSummary{
		Name:       c.Name,
		Operations: c.Operations,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}


