package eam

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// Service defines the boundary for enterprise asset state and maintenance.
type Service interface {
	GetAssetWithLogs(ctx context.Context, assetID uuid.UUID) (*Asset, []*MaintenanceLog, error)
	InsertAsset(ctx context.Context, asset *Asset) error
	InsertMaintenanceLog(ctx context.Context, log *MaintenanceLog) error
}

type eamService struct {
	client *spanner.Client
}

// NewService provisions the EAM service with the dedicated Spanner client.
func NewService(client *spanner.Client) Service {
	return &eamService{client: client}
}

// GetAssetWithLogs returns the physical asset and its interleaved maintenance history.
// It leverages a ReadOnlyTransaction snapshot to perfectly isolate the multi-table reads.
func (s *eamService) GetAssetWithLogs(ctx context.Context, assetID uuid.UUID) (*Asset, []*MaintenanceLog, error) {
	txn := s.client.ReadOnlyTransaction()
	defer txn.Close()

	// 1. Snapshot read of the Parent Asset Record
	row, err := txn.ReadRow(ctx, "Assets", spanner.Key{assetID.String()}, []string{
		"ID", "Name", "Type", "Status", "FinanceAssetID", "CreatedAt", "UpdatedAt",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read asset: %w", err)
	}

	var asset Asset
	if err := row.ToStruct(&asset); err != nil {
		return nil, nil, fmt.Errorf("failed to decode asset payload: %w", err)
	}

	// 2. Consistent snapshot query for interleaved MaintenanceLog data
	stmt := spanner.Statement{
		SQL: `SELECT AssetID, ID, TechnicianID, Description, Cost, CompletedAt, CreatedAt 
              FROM MaintenanceLogs 
              WHERE AssetID = @asset_id
              ORDER BY CompletedAt DESC`,
		Params: map[string]interface{}{
			"asset_id": assetID.String(),
		},
	}
	
	iter := txn.Query(ctx, stmt)
	defer iter.Stop()

	var logs []*MaintenanceLog
	for {
		logRow, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("maintenance log query traversal failed: %w", err)
		}
		var mLog MaintenanceLog
		if err := logRow.ToStruct(&mLog); err != nil {
			return nil, nil, fmt.Errorf("failed to decode log struct: %w", err)
		}
		logs = append(logs, &mLog)
	}

	return &asset, logs, nil
}

// InsertAsset tracks a new piece of corporate infrastructure.
func (s *eamService) InsertAsset(ctx context.Context, asset *Asset) error {
	mut, err := spanner.InsertStruct("Assets", asset)
	if err != nil {
		return err
	}
	_, err = s.client.Apply(ctx, []*spanner.Mutation{mut})
	return err
}

// InsertMaintenanceLog commits a maintenance operation bound into the Asset hierarchy.
func (s *eamService) InsertMaintenanceLog(ctx context.Context, log *MaintenanceLog) error {
	mut, err := spanner.InsertStruct("MaintenanceLogs", log)
	if err != nil {
		return err
	}
	_, err = s.client.Apply(ctx, []*spanner.Mutation{mut})
	return err
}
