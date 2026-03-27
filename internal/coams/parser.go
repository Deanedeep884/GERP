package coams

import (
	"bytes"
	"strings"

	"github.com/google/uuid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Parser struct holds the logic to process markdown into chunks and extract edges.
type Parser struct {
	md goldmark.Markdown
}

// NewParser initializes a new COAMS Markdown Parser
func NewParser() *Parser {
	return &Parser{
		md: goldmark.New(),
	}
}

// ParseResult holds the outcome of a markdown parsing operation
type ParseResult struct {
	Chunks []Chunk
	Edges  []Edge
}

// Parse processes the raw markdown content, slices it into semantic header chunks, and extracts all links.
func (p *Parser) Parse(channelID string, sourceDocID uuid.UUID, content []byte) (*ParseResult, error) {
	reader := text.NewReader(content)
	doc := p.md.Parser().Parse(reader)

	var chunks []Chunk
	var edges []Edge

	var currentHeaderPath []string
	var currentChunkBuf bytes.Buffer

	// FlushChunk is a helper to save the current buffer into a Chunk and reset it
	flushChunk := func() {
		if currentChunkBuf.Len() > 0 {
			headerPath := strings.Join(currentHeaderPath, " -> ")
			textStr := currentChunkBuf.String()
			chunks = append(chunks, Chunk{
				ID:         uuid.New(),
				ChannelID:  channelID,
				DocumentID: sourceDocID,
				HeaderPath: headerPath,
				Content:    textStr,
				Tokens:     len(strings.Fields(textStr)), // Naive token count proxy
			})
			currentChunkBuf.Reset()
		}
	}

	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		// We only process on entering the node for simplicity of extraction, except for blocks where we might want the text.
		
		switch node := n.(type) {
		case *ast.Heading:
			if entering {
				flushChunk() // Flush previous chunk before starting a new header
				
				// Update current header path based on level (H2, H3, etc)
				level := node.Level
				headingText := string(node.Text(content))
				
				// Keep path array scaled to the depth
				if len(currentHeaderPath) >= level {
					currentHeaderPath = currentHeaderPath[:level-1]
				}
				currentHeaderPath = append(currentHeaderPath, headingText)
				
				currentChunkBuf.WriteString(strings.Repeat("#", level) + " " + headingText + "\n\n")
			}
		case *ast.Link:
			if entering {
				url := string(node.Destination)
				edge := Edge{
					ID:               uuid.New(),
					ChannelID:        channelID,
					SourceDocumentID: sourceDocID,
				}

				if strings.HasPrefix(url, "doc:") {
					targetIDStr := strings.TrimPrefix(url, "doc:")
					targetUUID, err := uuid.Parse(targetIDStr)
					if err == nil {
						edge.TargetDocumentID = &targetUUID
						edge.IsExternal = false
					}
				} else {
					edge.IsExternal = true
					edge.ExternalURL = &url
				}
				edges = append(edges, edge)
			}
		case *ast.Text:
			if entering {
				currentChunkBuf.Write(node.Segment.Value(content))
			}
		}
		
		return ast.WalkContinue, nil
	})

	flushChunk() // Flush the final chunk
	
	return &ParseResult{
		Chunks: chunks,
		Edges:  edges,
	}, nil
}
