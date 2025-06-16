package qdrant

import (
	"context"
	"errors"
	"github.com/hl540/rag/embedding"
	"github.com/hl540/rag/vectorstore"
	"github.com/qdrant/go-client/qdrant"
)

type VectorStore struct {
	client   *qdrant.Client
	config   *qdrant.Config
	embedder embedding.Embedder
}

func New(opts ...Option) (vectorstore.VectorStore, error) {
	store := &VectorStore{config: &qdrant.Config{}}
	for _, opt := range opts {
		opt(store)
	}
	var err error
	store.client, err = qdrant.NewClient(store.config)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (v *VectorStore) createCollection(ctx context.Context, name string) error {
	exists, err := v.client.CollectionExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	embed, err := v.embedder.Embed(ctx, "test")
	if err != nil {
		return err
	}
	return v.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(len(embed)),
			Distance: qdrant.Distance_Dot,
		}),
	})
}

func (v *VectorStore) AddDocuments(ctx context.Context, name string, docs []*vectorstore.Document) error {
	if len(docs) == 0 {
		return errors.New("empty documents")
	}

	if err := v.createCollection(ctx, name); err != nil {
		return err
	}

	texts := make([]string, 0, len(docs))
	for _, doc := range docs {
		texts = append(texts, doc.Text)
	}
	embeds, err := v.embedder.Embeds(ctx, texts)
	if err != nil {
		return err
	}
	points := make([]*qdrant.PointStruct, 0, len(embeds))
	for i, doc := range docs {
		embed := embeds[i]
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewID(doc.Id),
			Vectors: qdrant.NewVectors(embed...),
			Payload: qdrant.NewValueMap(doc.Metadata),
		})
	}
	wait := true
	_, err = v.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: name,
		Wait:           &wait,
		Points:         points,
	})
	return err
}

func (v *VectorStore) SimilaritySearch(ctx context.Context, name string, query string, topK int) ([]*vectorstore.SearchResult, error) {
	if len(query) == 0 {
		return nil, errors.New("empty query")
	}

	embed, err := v.embedder.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	limit := uint64(topK)
	searchResult, err := v.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: name,
		Query:          qdrant.NewQuery(embed...),
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          &limit,
	})
	if err != nil {
		return nil, err
	}
	docs := make([]*vectorstore.SearchResult, 0, len(searchResult))
	for _, point := range searchResult {
		doc := &vectorstore.SearchResult{
			ID:       point.Id.GetUuid(),
			Score:    point.Score,
			Metadata: make(map[string]any),
		}
		for key, value := range point.Payload {
			doc.Metadata[key] = value.GetStringValue()
		}
		docs = append(docs, doc)
	}
	return docs, nil
}
