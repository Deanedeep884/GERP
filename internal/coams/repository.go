package coams

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the purely isolated interface over the partitioned AlloyDB storage.
// COAMS remains ignorant of connection pooling; it just expects this contract to be fulfilled.
type Repository interface {
	// Execute within an ACID transaction
	WithTransaction(ctx context.Context, fn func(txRepo Repository) error) error

	// Documents
	SaveDocument(ctx context.Context, doc *Document) error
	GetDocumentMetrics(ctx context.Context, channelID string, id uuid.UUID) (*Document, error)
	GetExistingDocumentIDs(ctx context.Context, channelID string, ids []uuid.UUID) (map[uuid.UUID]bool, error)

	// Chunks & Vectors
	SaveChunks(ctx context.Context, chunks []Chunk) error
	DeleteChunksByDocument(ctx context.Context, channelID string, documentID uuid.UUID) error
	SemanticSearch(ctx context.Context, channelID string, embedding []float32, limit int) ([]Chunk, error)

	// Links (Agent-Index)
	SaveEdges(ctx context.Context, edges []Edge) error
	DeleteEdgesByDocument(ctx context.Context, channelID string, documentID uuid.UUID) error
}
