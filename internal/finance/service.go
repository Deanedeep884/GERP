package finance

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the interface for the Finance domain resolvers and activities.
type Service interface {
	GetAccount(ctx context.Context, id uuid.UUID) (*Account, error)
	GetLedgerEntry(ctx context.Context, id uuid.UUID) (*LedgerEntry, []*LineItem, error)
	InsertLedgerEntry(ctx context.Context, entry *LedgerEntry, lines []*LineItem) error
}

type financeService struct {
	client *spanner.Client
}

// NewService safely instantiates the finance service with the Spanner client.
func NewService(client *spanner.Client) Service {
	return &financeService{client: client}
}

// GetAccount parses the Account by UUID natively.
func (s *financeService) GetAccount(ctx context.Context, id uuid.UUID) (*Account, error) {
	row, err := s.client.Single().ReadRow(ctx, "Accounts", spanner.Key{id.String()}, []string{
		"ID", "Name", "Type", "AccountOwnerID", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, fmt.Errorf("spanner read failed: %w", err)
	}

	var acc Account
	if err := row.ToStruct(&acc); err != nil {
		return nil, fmt.Errorf("failed to parse account struct: %w", err)
	}

	return &acc, nil
}

// GetLedgerEntry accurately reads an entry and mathematically proves interleaved LineItems 
// via a consistent read-only transaction snapshot to enforce ACID isolation.
func (s *financeService) GetLedgerEntry(ctx context.Context, id uuid.UUID) (*LedgerEntry, []*LineItem, error) {
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	row, err := txn.ReadRow(ctx, "LedgerEntries", spanner.Key{id.String()}, []string{
		"ID", "TransactionID", "Description", "CreatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ledger entry: %w", err)
	}

	var entry LedgerEntry
	if err := row.ToStruct(&entry); err != nil {
		return nil, nil, err
	}

	stmt := spanner.Statement{
		SQL: `SELECT LedgerEntryID, LineItemID, AccountID, AmountCents, CustomerID, CreatedAt 
              FROM LineItems 
              WHERE LedgerEntryID = @id`,
		Params: map[string]interface{}{
			"id": id.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var lines []*LineItem
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("line items query failed: %w", err)
		}
		var line LineItem
		if err := row.ToStruct(&line); err != nil {
			return nil, nil, err
		}
		lines = append(lines, &line)
	}

	return &entry, lines, nil
}

// InsertLedgerEntry commits a balanced double-entry transaction reliably via Spanner.
func (s *financeService) InsertLedgerEntry(ctx context.Context, entry *LedgerEntry, lines []*LineItem) error {
	var mutations []*spanner.Mutation

	// 1. Queue Parent Entry
	mut, err := spanner.InsertStruct("LedgerEntries", entry)
	if err != nil {
		return err
	}
	mutations = append(mutations, mut)

	// 2. Validate double-entry accounting balanced sum
	var totalCents int64
	for _, line := range lines {
		totalCents += line.AmountCents
		
		lineMut, err := spanner.InsertStruct("LineItems", line)
		if err != nil {
			return err
		}
		mutations = append(mutations, lineMut)
	}

	// Double-entry guarantee
	if totalCents != 0 {
		return fmt.Errorf("transaction not balanced (sum = %d cents)", totalCents)
	}

	// 3. Commit atomically to Spanner Hot Path
	_, err = s.client.Apply(ctx, mutations)
	return err
}
