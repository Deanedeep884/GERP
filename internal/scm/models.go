package scm

import (
	"time"

	"github.com/google/uuid"
)

// Product represents the master catalog definition of a physical good.
type Product struct {
	ID          uuid.UUID `json:"id" spanner:"ID"`
	SKU         string    `json:"sku" spanner:"SKU"`
	Name        string    `json:"name" spanner:"Name"`
	Description string    `json:"description" spanner:"Description"`
	IsActive    bool      `json:"is_active" spanner:"IsActive"`
	CreatedAt   time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt   time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// InventoryLot tracks actual physical instances of a Product in a Warehouse limit.
type InventoryLot struct {
	ID          uuid.UUID `json:"id" spanner:"ID"`
	ProductID   uuid.UUID `json:"product_id" spanner:"ProductID"`   // Parent Product
	WarehouseID uuid.UUID `json:"warehouse_id" spanner:"WarehouseID"` // Golden Thread mapping to `internal/eam`
	Quantity    int       `json:"quantity" spanner:"Quantity"`      // Count
	CostBasis   int64     `json:"cost_basis" spanner:"CostBasis"`   // Financial tracking (cents)
	CreatedAt   time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt   time.Time `json:"updated_at" spanner:"UpdatedAt"`
}
