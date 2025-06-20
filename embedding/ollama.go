package embedding

import (
	"context"
	"errors"
	"github.com/ollama/ollama/api"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// OllamaEmbedder 是一个基于 Ollama 的向量嵌入器
type OllamaEmbedder struct {
	client *api.Client
	model  string
	// 限制并发数量的信号量
	sem *semaphore.Weighted
}

// NewOllamaEmbedder 创建一个新的 Ollama OllamaEmbedder 实例
func NewOllamaEmbedder(client *api.Client, model string) Embedder {
	// 默认限制最大并发数为 5
	return &OllamaEmbedder{
		client: client,
		model:  model,
		sem:    semaphore.NewWeighted(5),
	}
}

// Embed 将单个文本转换为向量嵌入
func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, errors.New("empty text")
	}

	embed, err := e.client.Embed(ctx, &api.EmbedRequest{
		Model: e.model,
		Input: text,
	})
	if err != nil {
		return nil, err
	}
	if len(embed.Embeddings) == 0 {
		return nil, errors.New("no embeddings returned")
	}
	return embed.Embeddings[0], nil
}

// Embeds 将多个文本并发转换为向量嵌入，但限制并发数量
func (e *OllamaEmbedder) Embeds(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("empty texts")
	}

	embeds := make([][]float32, len(texts))
	g, ctx := errgroup.WithContext(ctx)

	for i, text := range texts {
		i, text := i, text // 创建新的变量以避免闭包问题
		if err := e.sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		g.Go(func() error {
			defer e.sem.Release(1)
			embed, err := e.Embed(ctx, text)
			if err != nil {
				return err
			}
			embeds[i] = embed
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return embeds, nil
}
