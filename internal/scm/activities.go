package scm

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

type InventoryRequest struct {
	LotID    uuid.UUID
	Quantity int       // Absolute requested count
	OrderID  uuid.UUID // Golden Thread context (Revenue domain)
}

// Activities groups the Temporal operations, isolating the stateless Spanner boundary.
type Activities struct {
	scmService Service
}

func NewActivities(svc Service) *Activities {
	return &Activities{
		scmService: svc,
	}
}

// AllocateInventoryActivity subtracts physical stock inside the supply chain bounds.
func (a *Activities) AllocateInventoryActivity(ctx context.Context, req InventoryRequest) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Attempting to allocate inventory", "LotID", req.LotID.String(), "OrderID", req.OrderID.String())

	if req.Quantity <= 0 {
		return fmt.Errorf("allocation quantity magnitude must be strictly positive")
	}

	// Committing resource block: requires a negative delta
	err := a.scmService.MutateInventory(ctx, req.LotID, -req.Quantity)
	if err != nil {
		logger.Error("Allocation failed", "Error", err)
		return fmt.Errorf("failed to allocate stock (potential OOS condition): %w", err)
	}

	logger.Info("Successfully allocated atomic stock bounds", "Quantity", req.Quantity)
	return nil
}

// ReverseInventoryActivity acts as the compensating invariant during a multi-domain saga rollback.
func (a *Activities) ReverseInventoryActivity(ctx context.Context, req InventoryRequest) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Reversing inventory allocation (Saga Rollback invoked)", "LotID", req.LotID.String(), "OrderID", req.OrderID.String())

	if req.Quantity <= 0 {
		return fmt.Errorf("rollback quantity magnitude must be strictly positive")
	}

	// Re-infusing resource block: requires additive delta
	err := a.scmService.MutateInventory(ctx, req.LotID, req.Quantity)
	if err != nil {
		logger.Error("Rollback restoration failed (CRITICAL PATH)", "Error", err)
		return fmt.Errorf("critical saga invariant broken: failed to restore locked stock: %w", err)
	}

	logger.Info("Successfully reversed and restored atomic stock allocation", "Quantity", req.Quantity)
	return nil
}
