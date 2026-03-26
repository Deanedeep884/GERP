package finance

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a financial ledger account base.
type Account struct {
	ID             uuid.UUID `json:"id" spanner:"ID"`
	Name           string    `json:"name" spanner:"Name"`
	Type           string    `json:"type" spanner:"Type"`               // e.g., "ASSET", "LIABILITY", "EQUITY", "REVENUE"
	AccountOwnerID uuid.UUID `json:"account_owner_id" spanner:"AccountOwnerID"` // Soft link to MDM or HCM via the Golden Thread
	CreatedAt      time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt      time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// LedgerEntry represents a single transaction event header.
type LedgerEntry struct {
	ID            uuid.UUID `json:"id" spanner:"ID"`
	TransactionID uuid.UUID `json:"transaction_id" spanner:"TransactionID"` // Soft link to orchestrating saga in pipeline
	Description   string    `json:"description" spanner:"Description"`
	CreatedAt     time.Time `json:"created_at" spanner:"CreatedAt"`
}

// LineItem represents the double-entry accounting records (Debits/Credits) belonging to a LedgerEntry.
type LineItem struct {
	LedgerEntryID uuid.UUID `json:"ledger_entry_id" spanner:"LedgerEntryID"` // Parent Key
	LineItemID    uuid.UUID `json:"line_item_id" spanner:"LineItemID"`     
	AccountID     uuid.UUID `json:"account_id" spanner:"AccountID"`
	AmountCents   int64     `json:"amount_cents" spanner:"AmountCents"` // Positive for Debit, Negative for Credit
	CustomerID    uuid.UUID `json:"customer_id" spanner:"CustomerID"`   // Soft link to Revenue/MDM profiles
	CreatedAt     time.Time `json:"created_at" spanner:"CreatedAt"`
}
