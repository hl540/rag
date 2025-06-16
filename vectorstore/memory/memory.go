package memory

import (
	"context"
	"errors"
	"github.com/hl540/rag/embedding"
	"github.com/hl540/rag/vectorstore"
	"math"
	"sort"
)

type VectorRecord struct {
	Id        string
	Text      string
	Embedding []float32
	Metadata  map[string]any
}

type VectorStore struct {
	embedder embedding.Embedder
	store    map[string][]*VectorRecord
}

func New(embedder embedding.Embedder) *VectorStore {
	return &VectorStore{
		embedder: embedder,
		store:    make(map[string][]*VectorRecord),
	}
}

func (v *VectorStore) AddDocuments(ctx context.Context, name string, docs []*vectorstore.Document) error {
	if len(docs) == 0 {
		return errors.New("empty documents")
	}
	if v.store[name] == nil {
		v.store[name] = make([]*VectorRecord, 0)
	}

	texts := make([]string, 0, len(docs))
	for _, doc := range docs {
		texts = append(texts, doc.Text)
	}
	embeds, err := v.embedder.Embeds(ctx, texts)
	if err != nil {
		return err
	}

	for i, doc := range docs {
		embed := embeds[i]
		v.store[name] = append(v.store[name], &VectorRecord{
			Id:        doc.Id,
			Text:      doc.Text,
			Embedding: embed,
			Metadata:  doc.Metadata,
		})
	}
	return nil
}

func (v *VectorStore) SimilaritySearch(ctx context.Context, name string, query string, topK int) ([]*vectorstore.SearchResult, error) {
	if len(query) == 0 {
		return nil, errors.New("empty query")
	}

	if v.store[name] == nil {
		return nil, errors.New("no such document")
	}

	embed, err := v.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	similarities := make([]*vectorstore.SearchResult, 0)
	for _, doc := range v.store[name] {
		similarity, err := v.CosineSimilarity(embed, doc.Embedding)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, &vectorstore.SearchResult{
			ID:       doc.Id,
			Score:    similarity,
			Metadata: doc.Metadata,
		})
	}
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})
	return similarities[:topK], nil
}

func (v *VectorStore) CosineSimilarity(vec1, vec2 []float32) (float32, error) {
	if len(vec1) != len(vec1) {
		return 0, errors.New("vectors must have the same length")
	}

	// Calculate dot product
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for i := 0; i < len(vec1); i++ {
		dotProduct += float64(vec1[i] * vec2[i])
		magnitude1 += float64(vec1[i] * vec1[i])
		magnitude2 += float64(vec2[i] * vec2[i])
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	// Check for zero magnitudes to avoid division by zero
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, errors.New("vector magnitude cannot be zero")
	}

	return float32(dotProduct / (magnitude1 * magnitude2)), nil
}
