package vectorstore

import "context"

type SearchResult struct {
	ID       string
	Score    float32
	Metadata map[string]any
}

type Document struct {
	Id       string
	Text     string
	Metadata map[string]any
}

const ContentKey = "content"

type VectorStore interface {
	AddDocument(ctx context.Context, name string, doc *Document) error
	AddDocuments(ctx context.Context, name string, docs []*Document) error
	SimilaritySearch(ctx context.Context, name string, query string, topK int) ([]*SearchResult, error)
}
