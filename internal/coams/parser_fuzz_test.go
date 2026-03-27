package coams_test

import (
	"testing"
	"gerp/internal/coams"
	"github.com/google/uuid"
)

func FuzzMarkdownASTParser(f *testing.F) {
	// Seed the fuzzer with expected inputs
	f.Add([]byte("# Standard Header\n[Valid Link](doc:123)"))
	f.Add([]byte("## Header with no text"))
	f.Add([]byte("[[[[Circular Broken Links]](doc:foo)"))

	f.Fuzz(func(t *testing.T, rawMarkdown []byte) {
		// The parser must NEVER panic, regardless of the input.
		// It should gracefully return an error or a safe, empty AST.
		parser := coams.NewParser()
		result, err := parser.Parse("engineering", uuid.New(), rawMarkdown)
		
		if err == nil && result != nil {
			// If it somehow parsed garbage without error, verify the memory bounds
			if len(result.Edges) > 10000 {
				t.Errorf("Parser allowed a massive edge array, potential memory exhaustion attack.")
			}
		}
	})
}
