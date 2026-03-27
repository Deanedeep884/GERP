package coams

import (
	"context"
)

// Embedder defines the contract for generating vector representations of semantic text chunks.
// During the Publish Saga Lifecycle, GERP will map this to Vertex AI's `textembedding-gecko`
// via its own internal IAM bindings and inject it into the COAMS ignorant engine.
type Embedder interface {
	// GenerateEmbeddings converts a batch of text strings into an array of float32 vectors.
	// E.g., calling Vertex AI through the official GCP Go SDK.
	GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}
