package scm

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the boundary for supply chain operations.
type Service interface {
	GetProductWithInventory(ctx context.Context, productID uuid.UUID) (*Product, []*InventoryLot, error)
	MutateInventory(ctx context.Context, lotID uuid.UUID, delta int) error
}

type scmService struct {
	client *spanner.Client
}

// NewService provisions the SCM service with the dedicated Spanner client.
func NewService(client *spanner.Client) Service {
	return &scmService{client: client}
}

// GetProductWithInventory uses a ReadOnlyTransaction to return a snapshot-consistent 
// aggregation of a product and its physical distributions.
func (s *scmService) GetProductWithInventory(ctx context.Context, productID uuid.UUID) (*Product, []*InventoryLot, error) {
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	// 1. Snapshot read of the Product record
	row, err := txn.ReadRow(ctx, "Products", spanner.Key{productID.String()}, []string{
		"ID", "SKU", "Name", "Description", "IsActive", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read product: %w", err)
	}

	var prod Product
	if err := row.ToStruct(&prod); err != nil {
		return nil, nil, fmt.Errorf("failed to decode product payload: %w", err)
	}

	// 2. Consistent query utilizing the secondary index
	stmt := spanner.Statement{
		SQL: `SELECT ID, ProductID, WarehouseID, Quantity, CostBasis, CreatedAt, UpdatedAt 
              FROM InventoryLots@{FORCE_INDEX=InventoryByProduct} 
              WHERE ProductID = @product_id`,
		Params: map[string]interface{}{
			"product_id": productID.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var lots []*InventoryLot
	for {
		lotRow, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("inventory query traversal failed: %w", err)
		}
		var lot InventoryLot
		if err := lotRow.ToStruct(&lot); err != nil {
			return nil, nil, fmt.Errorf("failed to decode lot struct: %w", err)
		}
		lots = append(lots, &lot)
	}

	return &prod, lots, nil
}

// MutateInventory performs an ACID read-modify-write spanning block to safely alter physical quantities.
func (s *scmService) MutateInventory(ctx context.Context, lotID uuid.UUID, delta int) error {
	_, err := s.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		row, err := txn.ReadRow(ctx, "InventoryLots", spanner.Key{lotID.String()}, []string{"Quantity"})
		if err != nil {
			return fmt.Errorf("failed to read lot quantity: %w", err)
		}

		var currentQuantity int64
		if err := row.ColumnByName("Quantity", &currentQuantity); err != nil {
			return fmt.Errorf("failed to parse existing quantity: %w", err)
		}

		newQuantity := currentQuantity + int64(delta)
		if newQuantity < 0 {
			// Mathematical boundary: Inventory cannot fall below true zero
			return fmt.Errorf("mutation rejected: insufficient physical stock (current=%d, required=%d)", currentQuantity, -delta)
		}

		mut := spanner.UpdateMap("InventoryLots", map[string]interface{}{
			"ID":        lotID.String(),
			"Quantity":  newQuantity,
			"UpdatedAt": spanner.CommitTimestamp,
		})

		return txn.BufferWrite([]*spanner.Mutation{mut})
	})

	return err
}
