package mdm

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the boundary for global identity resolution and master data mapping.
type Service interface {
	GetGlobalEntityWithMappings(ctx context.Context, id uuid.UUID) (*GlobalEntity, []*EntityMapping, error)
}

type mdmService struct {
	client *spanner.Client
}

// NewService provisions the MDM service with the dedicated Spanner client.
func NewService(client *spanner.Client) Service {
	return &mdmService{client: client}
}

// GetGlobalEntityWithMappings retrieves the Golden Record and resolves all domain-specific 
// identities in a single highly-optimized snapshot read block.
func (s *mdmService) GetGlobalEntityWithMappings(ctx context.Context, id uuid.UUID) (*GlobalEntity, []*EntityMapping, error) {
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	// 1. Snapshot read of the True Parent Reference
	row, err := txn.ReadRow(ctx, "GlobalEntities", spanner.Key{id.String()}, []string{
		"ID", "LegalName", "TaxID", "CountryCode", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read global entity: %w", err)
	}

	var entity GlobalEntity
	if err := row.ToStruct(&entity); err != nil {
		return nil, nil, fmt.Errorf("failed to decode global entity payload: %w", err)
	}

	// 2. Consistent query utilizing the Spanner interleave parent scope
	stmt := spanner.Statement{
		SQL: `SELECT GlobalEntityID, Domain, LocalID, CreatedAt, UpdatedAt 
              FROM EntityMappings 
              WHERE GlobalEntityID = @global_entity_id`,
		Params: map[string]interface{}{
			"global_entity_id": id.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var mappings []*EntityMapping
	for {
		mRow, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("mapping queries failed: %w", err)
		}
		var mp EntityMapping
		if err := mRow.ToStruct(&mp); err != nil {
			return nil, nil, fmt.Errorf("failed to decode entity mapping struct: %w", err)
		}
		mappings = append(mappings, &mp)
	}

	return &entity, mappings, nil
}
