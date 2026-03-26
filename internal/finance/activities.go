package finance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
)

// ChargeLedgerRequest defines the input payload for the Temporal execution vector.
type ChargeLedgerRequest struct {
	TransactionID   uuid.UUID
	AccountDebitID  uuid.UUID
	AccountCreditID uuid.UUID
	AmountCents     int64
	CustomerID      uuid.UUID // Connects the Golden Thread to MDM natively
	Description     string
}

// Activities groups the Temporal activity boundaries for the finance module,
// isolating the core Spanner Service interface.
type Activities struct {
	financeService Service
}

// NewActivities returns an instantiated Activities handler.
func NewActivities(svc Service) *Activities {
	return &Activities{
		financeService: svc,
	}
}

// ChargeLedgerActivity is a stateless Temporal Activity running on the distributed workers.
// It intercepts the cross-domain orchestrator payload and natively commits the database lock.
func (a *Activities) ChargeLedgerActivity(ctx context.Context, req ChargeLedgerRequest) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Attempting to lock Spanner Ledger", "TransactionID", req.TransactionID.String())

	ledgerID := uuid.New()

	entry := &LedgerEntry{
		ID:            ledgerID,
		TransactionID: req.TransactionID,
		Description:   req.Description,
		CreatedAt:     time.Now().UTC(), 
	}

	// Build the balanced entries
	debit := &LineItem{
		LedgerEntryID: ledgerID,
		LineItemID:    uuid.New(),
		AccountID:     req.AccountDebitID,
		AmountCents:   req.AmountCents, // Pos Debit
		CustomerID:    req.CustomerID,
		CreatedAt:     time.Now().UTC(),
	}

	credit := &LineItem{
		LedgerEntryID: ledgerID,
		LineItemID:    uuid.New(),
		AccountID:     req.AccountCreditID,
		AmountCents:   -req.AmountCents, // Neg Credit to balance
		CustomerID:    req.CustomerID,
		CreatedAt:     time.Now().UTC(),
	}

	// Invoke the strict domain integration interface
	err := a.financeService.InsertLedgerEntry(ctx, entry, []*LineItem{debit, credit})
	if err != nil {
		logger.Error("Failed to apply ledger transaction", "Error", err)
		return fmt.Errorf("spanner double-entry lock failed: %w", err)
	}

	logger.Info("Successfully bound cross-domain transaction in Spanner", "LedgerID", ledgerID.String())
	return nil
}
