package coams

import (
	"time"

	"github.com/google/uuid"
)

// Document represents the primary markdown entity within a specific channel partition.
type Document struct {
	ID          uuid.UUID
	TenantID    string
	ChannelID   string
	Title       string
	RawMarkdown string
	
	// Verbose Metadata mapped to IAM workspace identities
	CreatedBy string
	UpdatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// Chunk represents the AI's semantic unit of knowledge attached to a Document.
type Chunk struct {
	ID         uuid.UUID
	ChannelID  string
	DocumentID uuid.UUID
	HeaderPath string
	Content    string
	Tokens     int
	Embedding  []float32
}

// Edge represents an outbound link from a document, forming the Agent-Index.
type Edge struct {
	ID               uuid.UUID
	ChannelID        string
	SourceDocumentID uuid.UUID
	TargetDocumentID *uuid.UUID // Nullable if external
	IsExternal       bool
	ExternalURL      *string    // Nullable if internal
}

// SchemaDefinition is passed to the orchestrator (GERP) to build GraphQL dynamically
type SchemaDefinition struct {
	ModelID   string
	ChannelID string
	Name      string
	Fields    []FieldDefinition
	Relations []EdgeDefinition
}

type FieldDefinition struct {
	Name string
	Type string
}

type EdgeDefinition struct {
	TargetModel string
	EdgeName    string
}
