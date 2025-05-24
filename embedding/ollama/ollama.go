package ollama

import (
	"context"
	"errors"
	"github.com/hl540/rag/embedding"
	"github.com/ollama/ollama/api"
)

type Embedder struct {
	client *api.Client
	model  string
}

func New(client *api.Client, model string) embedding.Embedder {
	return &Embedder{
		client: client,
		model:  model,
	}
}

func (e *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
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
		return nil, errors.New("no embeddings")
	}
	return embed.Embeddings[0], nil
}

func (e *Embedder) Embeds(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("empty text")
	}
	embeds := make([][]float32, 0, len(texts))
	for _, text := range texts {
		embed, err := e.Embed(ctx, text)
		if err != nil {
			return nil, err
		}
		embeds = append(embeds, embed)
	}
	return embeds, nil
}
