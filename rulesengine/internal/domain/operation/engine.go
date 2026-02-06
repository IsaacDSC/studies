package operation

import (
	"context"
	"encoding/json"
	"sync"

	zen "github.com/gorules/zen-go"

	"rulesengine/internal/domain/blocklist"
	"rulesengine/internal/domain/campaigns"
	"rulesengine/internal/shared"
)

// Engine handles rule evaluation
type Engine struct {
	zenEngine     zen.Engine
	campaignsRepo campaigns.Repository
	blocklistSvc  *blocklist.Service
}

// NewEngine creates a new operation engine
func NewEngine(campaignsRepo campaigns.Repository, blocklistSvc *blocklist.Service) *Engine {
	return &Engine{
		zenEngine:     zen.NewEngine(zen.EngineConfig{}),
		campaignsRepo: campaignsRepo,
		blocklistSvc:  blocklistSvc,
	}
}

// Evaluate evaluates rules for the given operation and input
func (e *Engine) Evaluate(ctx context.Context, operation string, input EvaluationInput) (*EvaluationResult, error) {
	// 1. Check blocklist
	blocked, err := e.blocklistSvc.IsBlocked(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return NewBlockedResult(input.UserID.String(), input.RequestedAmount), nil
	}

	// 2. Load campaigns for operation
	campaignList, err := e.campaignsRepo.ListByOperation(ctx, operation)
	if err != nil {
		return nil, err
	}
	if len(campaignList) == 0 {
		return nil, shared.ErrNoCampaigns
	}

	// 3. Execute rules in parallel
	offers := e.executeParallel(ctx, campaignList, input)

	// 4. Aggregate results
	return e.aggregate(input, offers), nil
}

// executeParallel executes all campaign rules in parallel
func (e *Engine) executeParallel(ctx context.Context, campaignList []*campaigns.Campaign, input EvaluationInput) []Offer {
	var wg sync.WaitGroup
	results := make(chan Offer, len(campaignList))

	for _, campaign := range campaignList {
		wg.Add(1)
		go func(c *campaigns.Campaign) {
			defer wg.Done()

			offer, err := e.evaluateCampaign(c, input)
			if err != nil {
				return // Skip failed evaluations
			}

			results <- *offer
		}(campaign)
	}

	// Close channel when all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var offers []Offer
	for offer := range results {
		offers = append(offers, offer)
	}
	return offers
}

// evaluateCampaign evaluates a single campaign
func (e *Engine) evaluateCampaign(campaign *campaigns.Campaign, input EvaluationInput) (*Offer, error) {
	decision, err := e.zenEngine.CreateDecision(campaign.Graph)
	if err != nil {
		return nil, err
	}

	evalInput := map[string]any{
		"requested_amount": input.RequestedAmount,
		"user_id":          input.UserID.String(),
	}

	result, err := decision.Evaluate(evalInput)
	if err != nil {
		return nil, err
	}

	// Parse result
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Result map[string]any `json:"result"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return e.parseOffer(campaign.Name, resp.Result), nil
}

// parseOffer parses the evaluation result into an Offer
func (e *Engine) parseOffer(campaignName string, result map[string]any) *Offer {
	offer := &Offer{
		Campaign:     campaignName,
		Approved:     false,
		Rate:         0,
		Installments: []int{},
		Value:        0,
	}

	if result == nil {
		return offer
	}

	// Parse approved
	if approved, ok := result["approved"].(bool); ok {
		offer.Approved = approved
	}

	// Parse rate
	if rate, ok := result["rate"].(float64); ok {
		offer.Rate = rate
	}

	// Parse installments
	if installments, ok := result["installments"].([]any); ok {
		for _, i := range installments {
			if v, ok := i.(float64); ok {
				offer.Installments = append(offer.Installments, int(v))
			}
		}
	}

	// Parse value
	if value, ok := result["value"].(float64); ok {
		offer.Value = int(value)
	}

	return offer
}

// aggregate aggregates offers into a final result
func (e *Engine) aggregate(input EvaluationInput, offers []Offer) *EvaluationResult {
	// Check if any offer is approved
	hasApproved := false
	for _, offer := range offers {
		if offer.Approved {
			hasApproved = true
			break
		}
	}

	if hasApproved {
		return NewApprovedResult(input.UserID.String(), input.RequestedAmount, offers)
	}

	return NewDeniedResult(
		input.UserID.String(),
		input.RequestedAmount,
		"credit_limit_exceeded",
		"The requested amount exceeds the available limit",
	)
}


