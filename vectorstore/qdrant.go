package vectorstore

import (
	"context"
	"errors"
	"github.com/hl540/rag/embedding"
	"github.com/qdrant/go-client/qdrant"
)

// QdrantStore 是一个基于 Qdrant 的向量存储，支持文档的添加和相似度搜索
type QdrantStore struct {
	client   *qdrant.Client
	config   *qdrant.Config
	embedder embedding.Embedder
}

// NewQdrantStore 创建一个新的 QdrantVectorStore 实例
func NewQdrantStore(opts ...QdrantOption) (VectorStore, error) {
	store := &QdrantStore{config: &qdrant.Config{}}
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

// createCollection 检查并创建 Qdrant 集合
func (v *QdrantStore) createCollection(ctx context.Context, name string) error {
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

// AddDocuments 将文档添加到 Qdrant 集合中，并生成向量嵌入
func (v *QdrantStore) AddDocuments(ctx context.Context, name string, docs []*Document) error {
	if len(docs) == 0 {
		return errors.New("empty documents")
	}

	if err := v.createCollection(ctx, name); err != nil {
		return err
	}

	// 设置批处理大小
	batchSize := 100
	wait := true

	// 分批处理文档
	for i := 0; i < len(docs); i += batchSize {
		end := i + batchSize
		if end > len(docs) {
			end = len(docs)
		}

		// 获取当前批次的文档
		batchDocs := docs[i:end]
		texts := make([]string, 0, len(batchDocs))
		for _, doc := range batchDocs {
			texts = append(texts, doc.Text)
		}

		// 生成当前批次的向量嵌入
		embeds, err := v.embedder.Embeds(ctx, texts)
		if err != nil {
			return err
		}

		// 构建当前批次的点
		points := make([]*qdrant.PointStruct, 0, len(embeds))
		for j, doc := range batchDocs {
			embed := embeds[j]
			points = append(points, &qdrant.PointStruct{
				Id:      qdrant.NewID(doc.Id),
				Vectors: qdrant.NewVectors(embed...),
				Payload: qdrant.NewValueMap(doc.Metadata),
			})
		}

		// 推送当前批次到 Qdrant
		_, err = v.client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: name,
			Wait:           &wait,
			Points:         points,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// SimilaritySearch 在 Qdrant 集合中搜索与查询最相似的 topK 条记录
func (v *QdrantStore) SimilaritySearch(ctx context.Context, name string, query string, topK int) ([]*SearchResult, error) {
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
	docs := make([]*SearchResult, 0, len(searchResult))
	for _, point := range searchResult {
		doc := &SearchResult{
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
