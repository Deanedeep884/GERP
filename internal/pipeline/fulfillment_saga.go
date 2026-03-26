package pipeline

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"gerp/internal/finance"
	"gerp/internal/scm"
)

// GlobalFulfillmentSaga orchestrates the master fulfillment pipeline routing.
// It executes mathematically strict SCM allocations followed by Finance ledger charging,
// guaranteeing compensating reversing logic dynamically upon physical transaction failures.
func GlobalFulfillmentSaga(ctx workflow.Context, req FulfillmentRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Initializing Global Fulfillment Saga bounds", "TransactionID", req.TransactionID.String())

	// Highly durable retry bounds optimized for the Cloud Spanner hot-path tolerances
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 2,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Keep an immutable slice of locally allocated physical inventory tracks
	var successfullyAllocated []scm.InventoryRequest

	// STEP 1: Execute SCM Allocations globally
	for _, item := range req.Items {
		inventoryReq := scm.InventoryRequest{
			LotID:    item.LotID,
			Quantity: item.Quantity,
			OrderID:  req.TransactionID, 
		}

		err := workflow.ExecuteActivity(ctx, "AllocateInventoryActivity", inventoryReq).Get(ctx, nil)
		if err != nil {
			logger.Error("Encountered hard failure allocating physical inventory", "LotID", item.LotID.String(), "Error", err)
			
			// Revert the actively pinned quantities dynamically
			errRevert := executeCompensatingRollbacks(ctx, successfullyAllocated)
			if errRevert != nil {
				logger.Error("CATASTROPHIC MATRIX FAILURE: Unable to restore locked supply chain items", "Error", errRevert)
			}
			return err
		}

		// Save the local state marker for standard rolling compensations
		successfullyAllocated = append(successfullyAllocated, inventoryReq)
	}

	logger.Info("Supply Chain allocation matrix successfully bound.")

	// STEP 2: Charge the Global Ledger structurally via internal/finance
	financeReq := finance.ChargeLedgerRequest{
		TransactionID:   req.TransactionID,
		AccountDebitID:  req.AccountDebitID,
		AccountCreditID: req.AccountCreditID,
		AmountCents:     req.TotalAmountCents,
		CustomerID:      req.CustomerID,
		Description:     "Automated Global Fulfillment Saga Charge",
	}

	err := workflow.ExecuteActivity(ctx, "ChargeLedgerActivity", financeReq).Get(ctx, nil)
	if err != nil {
		logger.Error("Finance domain rejected ledger matrix block", "TransactionID", req.TransactionID.String(), "Error", err)

		// STEP 3: The Rollback - reverse the temporal memory graph of physical stock
		errRevert := executeCompensatingRollbacks(ctx, successfullyAllocated)
		if errRevert != nil {
			logger.Error("CATASTROPHIC MATRIX FAILURE: Ledger rejected and SCM rollbacks aborted", "Error", errRevert)
		}
		
		// Return the parent error structure enforcing saga registration failure
		return err
	}

	logger.Info("Global Fulfillment Saga correctly orchestrated and successfully finalized across all matrices.")
	return nil
}

// executeCompensatingRollbacks safely walks backwards physically over the verified tracks in the matrix bounds.
func executeCompensatingRollbacks(ctx workflow.Context, allocatedLots []scm.InventoryRequest) error {
	logger := workflow.GetLogger(ctx)

	for _, lot := range allocatedLots {
		logger.Info("Attempting compensating transaction rollback", "LotID", lot.LotID.String())
		
		// Wait mathematically for the execution loop return
		err := workflow.ExecuteActivity(ctx, "ReverseInventoryActivity", lot).Get(ctx, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
