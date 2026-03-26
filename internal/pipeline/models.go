package pipeline

import (
	"github.com/google/uuid"
)

// OrderItem specifies a discrete SKU quantity requirement within the matrix.
type OrderItem struct {
	LotID    uuid.UUID `json:"lot_id"`    // Golden Thread traversing to SCM InventoryLot
	Quantity int       `json:"quantity"`  // Physical discrete count mapping
}

// FulfillmentRequest defines the strict entry payload for the Temporal executor saga.
type FulfillmentRequest struct {
	TransactionID    uuid.UUID     `json:"transaction_id"`
	CustomerID       uuid.UUID     `json:"customer_id"`       // Golden Thread traversal constraint
	AccountDebitID   uuid.UUID     `json:"account_debit_id"`  // Finance Domain Target
	AccountCreditID  uuid.UUID     `json:"account_credit_id"` // Finance Domain Target
	TotalAmountCents int64         `json:"total_amount_cents"`
	Items            []OrderItem   `json:"items"`
}
