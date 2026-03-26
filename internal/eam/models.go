package eam

import (
	"time"

	"github.com/google/uuid"
)

// Asset represents a physically tracked asset managed by the enterprise (warehouses, delivery trucks, IoT).
type Asset struct {
	ID             uuid.UUID `json:"id" spanner:"ID"`
	Name           string    `json:"name" spanner:"Name"`
	Type           string    `json:"type" spanner:"Type"`             // e.g., "WAREHOUSE", "TRUCK", "HVAC"
	Status         string    `json:"status" spanner:"Status"`         // e.g., "ONLINE", "MAINTENANCE", "DECOMMISSIONED"
	FinanceAssetID uuid.UUID `json:"finance_asset_id" spanner:"FinanceAssetID"` // Golden Thread: Fixed Asset ledger pointer
	CreatedAt      time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt      time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// MaintenanceLog tracks the work history interleaved directly under a physical asset.
type MaintenanceLog struct {
	AssetID      uuid.UUID `json:"asset_id" spanner:"AssetID"`           // Parent Key
	ID           uuid.UUID `json:"id" spanner:"ID"`
	TechnicianID uuid.UUID `json:"technician_id" spanner:"TechnicianID"` // Golden Thread: HCM Employee mapping
	Description  string    `json:"description" spanner:"Description"`
	Cost         int64     `json:"cost" spanner:"Cost"`                  // Minor units (cents)
	CompletedAt  time.Time `json:"completed_at" spanner:"CompletedAt"`
	CreatedAt    time.Time `json:"created_at" spanner:"CreatedAt"`
}
