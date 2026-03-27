package coams

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Verifier ensures that a document's outbound semantic links (Agent-Index)
// are mathematically sound within the isolated channel partition context.
type Verifier struct {
	repo Repository
}

// NewVerifier initializes the Verifier with a provided repository.
func NewVerifier(repo Repository) *Verifier {
	return &Verifier{repo: repo}
}

// EnsureLinkIntegrity performs a set-difference verification between the requested outbound doc:uuid links
// and the actually existing records inside that specific channel's AlloyDB partition.
func (v *Verifier) EnsureLinkIntegrity(ctx context.Context, channelID string, edges []Edge) error {
	var targetIDs []uuid.UUID
	
	// Fast path loop to collect unique target UUIDs
	uniqueTargets := make(map[uuid.UUID]bool)
	for _, edge := range edges {
		if !edge.IsExternal && edge.TargetDocumentID != nil {
			uniqueTargets[*edge.TargetDocumentID] = true
		}
	}
	
	if len(uniqueTargets) == 0 {
		return nil // No internal graph links to mathematically verify
	}

	for id := range uniqueTargets {
		targetIDs = append(targetIDs, id)
	}

	// Fetch existing IDs from the repository's partitioned view
	existingSet, err := v.repo.GetExistingDocumentIDs(ctx, channelID, targetIDs)
	if err != nil {
		return fmt.Errorf("failed to query agent-index for verification: %w", err)
	}

	// Compute set difference
	var brokenLinks []uuid.UUID
	for reqID := range uniqueTargets {
		if !existingSet[reqID] {
			brokenLinks = append(brokenLinks, reqID)
		}
	}

	if len(brokenLinks) > 0 {
		return fmt.Errorf("mathematical link integrity failure in channel '%s': target documents %v do not exist or lack IAM visibility", channelID, brokenLinks)
	}

	return nil
}
