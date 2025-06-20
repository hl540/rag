package vectorstore

import (
	"context"
	"errors"
	"github.com/hl540/rag/embedding"
	"math"
	"sort"
)

// MemoryVectorRecord 表示一条向量记录，包含 ID、文本、向量嵌入和元数据
type MemoryVectorRecord struct {
	Id        string
	Text      string
	Embedding []float32
	Metadata  map[string]any
}

// MemoryStore 是一个内存向量存储，支持按名称分组存储向量记录
type MemoryStore struct {
	embedder embedding.Embedder
	store    map[string][]*MemoryVectorRecord
}

// NewMemoryStore 创建一个新的 VectorStore 实例
func NewMemoryStore(embedder embedding.Embedder) VectorStore {
	return &MemoryStore{
		embedder: embedder,
		store:    make(map[string][]*MemoryVectorRecord),
	}
}

// AddDocuments 将文档添加到指定名称的存储中，并生成向量嵌入
func (v *MemoryStore) AddDocuments(ctx context.Context, name string, docs []*Document) error {
	if len(docs) == 0 {
		return errors.New("empty documents")
	}
	if v.store[name] == nil {
		v.store[name] = make([]*MemoryVectorRecord, 0)
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
		v.store[name] = append(v.store[name], &MemoryVectorRecord{
			Id:        doc.Id,
			Text:      doc.Text,
			Embedding: embed,
			Metadata:  doc.Metadata,
		})
	}
	return nil
}

// SimilaritySearch 在指定名称的存储中搜索与查询最相似的 topK 条记录
func (v *MemoryStore) SimilaritySearch(ctx context.Context, name string, query string, topK int) ([]*SearchResult, error) {
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

	similarities := make([]*SearchResult, 0)
	for _, doc := range v.store[name] {
		similarity, err := v.CosineSimilarity(embed, doc.Embedding)
		if err != nil {
			return nil, err
		}
		similarities = append(similarities, &SearchResult{
			ID:       doc.Id,
			Score:    similarity,
			Metadata: doc.Metadata,
		})
	}
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Score > similarities[j].Score
	})
	if topK > len(similarities) {
		topK = len(similarities)
	}
	return similarities[:topK], nil
}

// CosineSimilarity 计算两个向量的余弦相似度
func (v *MemoryStore) CosineSimilarity(vec1, vec2 []float32) (float32, error) {
	if len(vec1) != len(vec2) {
		return 0, errors.New("vectors must have the same length")
	}

	// 计算点积和向量模长
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

	// 检查向量模长是否为零，避免除零错误
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0, errors.New("vector magnitude cannot be zero")
	}

	return float32(dotProduct / (magnitude1 * magnitude2)), nil
}
