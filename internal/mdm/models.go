package mdm

import (
	"time"

	"github.com/google/uuid"
)

// GlobalEntity serves as the absolute canonical representation of a real-world corporation or individual.
type GlobalEntity struct {
	ID          uuid.UUID `json:"id" spanner:"ID"`
	LegalName   string    `json:"legal_name" spanner:"LegalName"`
	TaxID       string    `json:"tax_id" spanner:"TaxID"`
	CountryCode string    `json:"country_code" spanner:"CountryCode"`
	CreatedAt   time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt   time.Time `json:"updated_at" spanner:"UpdatedAt"`
}

// EntityMapping translates a GlobalEntity down into domain-specific representations (e.g., 'Customer', 'Vendor').
type EntityMapping struct {
	GlobalEntityID uuid.UUID `json:"global_entity_id" spanner:"GlobalEntityID"` // Parent Key
	Domain         string    `json:"domain" spanner:"Domain"`                   // e.g., "revenue", "scm"
	LocalID        uuid.UUID `json:"local_id" spanner:"LocalID"`                // Golden Thread to domain-specific identity record
	CreatedAt      time.Time `json:"created_at" spanner:"CreatedAt"`
	UpdatedAt      time.Time `json:"updated_at" spanner:"UpdatedAt"`
}
